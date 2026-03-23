package ai_engine

import (
	"backend/internal/domain"
	pb "backend/internal/proto/ai_engine/v1"
	"backend/internal/temporal/activity"
	"backend/internal/temporal/workflow"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"go.temporal.io/sdk/client"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *Handler) GetCandidateScore(ctx context.Context, req *pb.GetCandidateScoreRequest) (*pb.GetCandidateScoreResponse, error) {
	score, factors, err := h.candidateRepo.GetScoreByCandidateID(ctx, req.GetCandidateId())
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &pb.GetCandidateScoreResponse{HasScore: false}, nil
		}
		return nil, status.Errorf(codes.Internal, "get candidate score: %v", err)
	}

	protoFactors := make([]*pb.ScoreFactor, len(factors))
	for i, f := range factors {
		protoFactors[i] = &pb.ScoreFactor{
			Type:        string(f.Type),
			Description: f.Description,
			Impact:      int32(f.Impact),
		}
	}

	return &pb.GetCandidateScoreResponse{
		HasScore:   true,
		MatchScore: int32(score.MatchScore),
		Factors:    protoFactors,
	}, nil
}

func (h *Handler) GetCandidateScores(ctx context.Context, req *pb.GetCandidateScoresRequest) (*pb.GetCandidateScoresResponse, error) {
	scores, err := h.candidateRepo.GetScoresByCandidateIDs(ctx, req.GetCandidateIds())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "get candidate scores: %v", err)
	}

	protoScores := make(map[string]*pb.GetCandidateScoreResponse)
	for id, score := range scores {
		protoScores[id] = &pb.GetCandidateScoreResponse{
			HasScore:   true,
			MatchScore: int32(score.MatchScore),
		}
	}

	return &pb.GetCandidateScoresResponse{
		Scores: protoScores,
	}, nil
}

func (h *Handler) GenerateCandidateEmail(ctx context.Context, req *pb.GenerateCandidateEmailRequest) (*pb.GenerateCandidateEmailResponse, error) {
	cand, _, _, err := h.candidateRepo.GetByID(ctx, req.GetCandidateId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "get candidate: %v", err)
	}

	var matchScore int
	score, _, scoreErr := h.candidateRepo.GetScoreByCandidateID(ctx, req.GetCandidateId())
	if scoreErr == nil {
		matchScore = score.MatchScore
	}

	var firstName, lastName, email string
	if cand.FirstName != nil {
		firstName = *cand.FirstName
	}
	if cand.LastName != nil {
		lastName = *cand.LastName
	}
	if cand.Email != nil {
		email = *cand.Email
	}

	input := activity.EmailGenerateInput{
		CandidateID:        cand.ID,
		GeneratedByUserID:  req.GetGeneratedByUserId(),
		CandidateFirstName: firstName,
		CandidateLastName:  lastName,
		CandidateEmail:     email,
		Role:               cand.JobTitle,
		Skills:             cand.Skills,
		MatchScore:         matchScore,
		TeamID:             req.GetTeamId(),
		EmailType:          mapEmailType(req.GetType()),
		Tone:               mapEmailTone(req.GetTone()),
		Locale:             req.GetLocale(),
		RecruiterName:      req.GetRecruiterName(),
		CompanyName:        req.GetCompanyName(),
	}

	workflowID := fmt.Sprintf("email-gen-%s-%s", req.GetCandidateId(), uuid.New().String())

	opts := client.StartWorkflowOptions{
		ID:                       workflowID,
		TaskQueue:                h.temporalClient.TaskQueue(),
		WorkflowExecutionTimeout: 3 * time.Minute,
	}

	we, err := h.temporalClient.TemporalClient.ExecuteWorkflow(ctx, opts, workflow.EmailGenerateWorkflow, input)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "start email workflow: %v", err)
	}

	var output activity.EmailGenerateOutput
	if err := we.Get(ctx, &output); err != nil {
		return nil, status.Errorf(codes.Internal, "email workflow failed: %v", err)
	}

	return &pb.GenerateCandidateEmailResponse{
		CommunicationId: output.CommunicationID,
		Subject:         output.Subject,
		Body:            output.Body,
	}, nil
}

func (h *Handler) GetCandidateEmails(ctx context.Context, req *pb.GetCandidateEmailsRequest) (*pb.GetCandidateEmailsResponse, error) {
	comms, err := h.commRepo.GetByCandidateID(ctx, req.GetCandidateId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "get candidate emails: %v", err)
	}

	entries := make([]*pb.CommunicationEntry, len(comms))
	for i, c := range comms {
		entries[i] = &pb.CommunicationEntry{
			CommunicationId:   c.ID,
			Type:              mapDomainEmailType(c.Type),
			Subject:           c.Subject,
			Body:              c.Body,
			GeneratedByUserId: c.GeneratedByUserID,
			IsSent:            c.SentAt != nil,
			GeneratedAtUnix:   c.CreatedAt.Unix(),
		}
	}

	return &pb.GetCandidateEmailsResponse{Emails: entries}, nil
}

func (h *Handler) CompareCandidates(ctx context.Context, req *pb.CompareCandidatesRequest) (*pb.CompareCandidatesResponse, error) {
	candidates := make([]activity.CandidateCompareCandidate, len(req.GetCandidates()))
	for i, c := range req.GetCandidates() {
		candidates[i] = activity.CandidateCompareCandidate{
			ID:         c.GetCandidateId(),
			FirstName:  c.GetFirstName(),
			LastName:   c.GetLastName(),
			Role:       c.GetRole(),
			MatchScore: int(c.GetMatchScore()),
			Skills:     c.GetSkills(),
			ParsedText: c.GetParsedText(),
		}
	}

	input := activity.CandidateCompareInput{
		TeamID:          req.GetTeamId(),
		Candidates:      candidates,
		JobRequirements: req.GetJobRequirements(),
		Locale:          req.GetLocale(),
	}

	workflowID := fmt.Sprintf("compare-%s", uuid.New().String())

	opts := client.StartWorkflowOptions{
		ID:                       workflowID,
		TaskQueue:                h.temporalClient.TaskQueue(),
		WorkflowExecutionTimeout: 3 * time.Minute,
	}

	we, err := h.temporalClient.TemporalClient.ExecuteWorkflow(ctx, opts, workflow.CompareCandidatesWorkflow, input)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "start compare workflow: %v", err)
	}

	var output activity.CandidateCompareOutput
	if err := we.Get(ctx, &output); err != nil {
		return nil, status.Errorf(codes.Internal, "compare workflow failed: %v", err)
	}

	results := make(map[string]*pb.CandidateComparisonEntry, len(output))
	for id, entry := range output {
		results[id] = &pb.CandidateComparisonEntry{
			Experience:  entry.Experience,
			Skills:      entry.Skills,
			SalaryRange: entry.SalaryRange,
			Risks:       entry.Risks,
			Summary:     entry.Summary,
		}
	}

	return &pb.CompareCandidatesResponse{Results: results}, nil
}

func (h *Handler) ParseJobDescription(ctx context.Context, req *pb.ParseJobDescriptionRequest) (*pb.ParseJobDescriptionResponse, error) {
	input := activity.JobParseInput{
		TeamID:  req.GetTeamId(),
		RawText: req.GetRawText(),
		Locale:  req.GetLocale(),
	}

	workflowID := fmt.Sprintf("job-parse-%s", uuid.New().String())

	opts := client.StartWorkflowOptions{
		ID:                       workflowID,
		TaskQueue:                h.temporalClient.TaskQueue(),
		WorkflowExecutionTimeout: 3 * time.Minute,
	}

	we, err := h.temporalClient.TemporalClient.ExecuteWorkflow(ctx, opts, workflow.JobParseWorkflow, input)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "start job parse workflow: %v", err)
	}

	var output activity.JobParseOutput
	if err := we.Get(ctx, &output); err != nil {
		return nil, status.Errorf(codes.Internal, "job parse workflow failed: %v", err)
	}

	return &pb.ParseJobDescriptionResponse{
		Title:        output.Title,
		Description:  output.Description,
		Requirements: output.Requirements,
		WorkFormat:   output.WorkFormat,
		SalaryMin:    int32(output.SalaryMin),
		SalaryMax:    int32(output.SalaryMax),
		Currency:     output.Currency,
	}, nil
}

func mapEmailType(t pb.EmailType) string {
	switch t {
	case pb.EmailType_EMAIL_TYPE_REJECTION:
		return string(domain.EmailRejection)
	case pb.EmailType_EMAIL_TYPE_INTERVIEW_INVITE:
		return string(domain.EmailInterviewInvite)
	default:
		return string(domain.EmailRejection)
	}
}

func mapEmailTone(t pb.EmailTone) string {
	switch t {
	case pb.EmailTone_EMAIL_TONE_PROFESSIONAL:
		return "professional"
	case pb.EmailTone_EMAIL_TONE_FRIENDLY:
		return "friendly"
	case pb.EmailTone_EMAIL_TONE_BRIEF:
		return "brief"
	default:
		return "professional"
	}
}

func mapDomainEmailType(t domain.EmailType) pb.EmailType {
	switch t {
	case domain.EmailRejection:
		return pb.EmailType_EMAIL_TYPE_REJECTION
	case domain.EmailInterviewInvite:
		return pb.EmailType_EMAIL_TYPE_INTERVIEW_INVITE
	default:
		return pb.EmailType_EMAIL_TYPE_UNSPECIFIED
	}
}
