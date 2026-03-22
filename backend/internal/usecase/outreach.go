package usecase

import (
	"backend/internal/domain"
	pbAI "backend/internal/proto/ai_engine/v1"
	"backend/internal/repo"
	"backend/pkg/email"
	"context"
	"fmt"
	"go.uber.org/zap"
)

type OutreachUseCase interface {
	GenerateEmail(ctx context.Context, candidateID, userID, role string, emailType domain.EmailType, tone string, locale string) (*domain.Communication, error)
	SendEmail(ctx context.Context, communicationID, subject, body string) error
}

type outreachUseCase struct {
	log           *zap.Logger
	candidateRepo repo.CandidateRepository
	commRepo      repo.CommunicationRepository
	aiUseCase     AiUseCase
	emailClient   *email.Client
}

func NewOutreachUseCase(
	log *zap.Logger,
	candidateRepo repo.CandidateRepository,
	commRepo repo.CommunicationRepository,
	aiUseCase AiUseCase,
	emailClient *email.Client,
) OutreachUseCase {
	return &outreachUseCase{
		log:           log,
		candidateRepo: candidateRepo,
		commRepo:      commRepo,
		aiUseCase:     aiUseCase,
		emailClient:   emailClient,
	}
}

func (u *outreachUseCase) GenerateEmail(ctx context.Context, candidateID, userID, role string, emailType domain.EmailType, tone string, locale string) (*domain.Communication, error) {
	pbType := pbAI.EmailType_EMAIL_TYPE_REJECTION
	if emailType == domain.EmailInterviewInvite {
		pbType = pbAI.EmailType_EMAIL_TYPE_INTERVIEW_INVITE
	}

	pbTone := pbAI.EmailTone_EMAIL_TONE_PROFESSIONAL
	switch tone {
	case "friendly":
		pbTone = pbAI.EmailTone_EMAIL_TONE_FRIENDLY
	case "brief":
		pbTone = pbAI.EmailTone_EMAIL_TONE_BRIEF
	}

	resp, err := u.aiUseCase.GenerateCandidateEmail(ctx, &pbAI.GenerateCandidateEmailRequest{
		CandidateId:       candidateID,
		GeneratedByUserId: userID,
		Type:               pbType,
		Tone:               pbTone,
		Locale:             locale,
	})
	if err != nil {
		return nil, fmt.Errorf("ai generate email: %w", err)
	}

	comm := &domain.Communication{
		ID:                resp.CommunicationId,
		CandidateID:       candidateID,
		GeneratedByUserID: userID,
		Type:              emailType,
		Subject:           resp.Subject,
		Body:              resp.Body,
	}

	return comm, nil
}

func (u *outreachUseCase) SendEmail(ctx context.Context, communicationID, subject, body string) error {
	comm, err := u.commRepo.GetByID(ctx, communicationID)
	if err != nil {
		return fmt.Errorf("get communication: %w", err)
	}

	if comm.SentAt != nil {
		return fmt.Errorf("email already sent")
	}

	cand, _, _, err := u.candidateRepo.GetByID(ctx, comm.CandidateID)
	if err != nil {
		return fmt.Errorf("get candidate: %w", err)
	}

	if cand.Email == nil || *cand.Email == "" {
		return fmt.Errorf("candidate email not found")
	}

	if err := u.emailClient.Send(ctx, *cand.Email, subject, body); err != nil {
		return fmt.Errorf("send email: %w", err)
	}

	if err := u.commRepo.MarkSent(ctx, communicationID); err != nil {
		u.log.Error("failed to mark communication as sent", zap.String("id", communicationID), zap.Error(err))
		// We don't return error here because the email was already sent.
	}

	return nil
}
