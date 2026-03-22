package usecase

import (
	"backend/internal/domain"
	pbAI "backend/internal/proto/ai_engine/v1"
	pbAuth "backend/internal/proto/auth/v1"
	"backend/internal/repo"
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
)

type DashboardUseCase interface {
	GetDashboardStats(ctx context.Context, teamID string) (*domain.DashboardStats, error)
	GetApplicationDynamics(ctx context.Context, req domain.DashboardDynamicsRequest, teamID string) ([]domain.ChartDataPoint, error)
	GetRecentActivity(ctx context.Context, teamID string, limit int) ([]domain.ActivityLogEntry, error)
}

type dashboardUseCase struct {
	log        *zap.Logger
	repo       repo.DashboardRepository
	authClient *pbAuth.AuthServiceClient
	aiClient   *pbAI.AIEngineServiceClient
}

func NewDashboardUseCase(
	log *zap.Logger,
	repo repo.DashboardRepository,
	authClient *pbAuth.AuthServiceClient,
	aiClient *pbAI.AIEngineServiceClient,
) DashboardUseCase {
	return &dashboardUseCase{
		log:        log,
		repo:       repo,
		authClient: authClient,
		aiClient:   aiClient,
	}
}

func (u *dashboardUseCase) GetDashboardStats(ctx context.Context, teamID string) (*domain.DashboardStats, error) {
	return u.repo.GetStats(ctx, teamID)
}

func (u *dashboardUseCase) GetApplicationDynamics(ctx context.Context, req domain.DashboardDynamicsRequest, teamID string) ([]domain.ChartDataPoint, error) {
	layout := "2006-01-02"
	start, err := time.Parse(layout, req.StartDate)
	if err != nil {
		return nil, fmt.Errorf("invalid start_date format (expected YYYY-MM-DD): %w", err)
	}
	end, err := time.Parse(layout, req.EndDate)
	if err != nil {
		return nil, fmt.Errorf("invalid end_date format (expected YYYY-MM-DD): %w", err)
	}

	if req.Type != "daily" && req.Type != "monthly" {
		return nil, fmt.Errorf("type must be daily or monthly")
	}

	return u.repo.GetApplicationDynamics(ctx, teamID, start, end, req.Type)
}

func (u *dashboardUseCase) GetRecentActivity(ctx context.Context, teamID string, limit int) ([]domain.ActivityLogEntry, error) {
	logs, err := u.repo.GetRecentActivity(ctx, teamID, limit)
	if err != nil {
		return nil, fmt.Errorf("get recent activity: %w", err)
	}

	if len(logs) == 0 {
		return logs, nil
	}

	var userIDs []string
	var candidateIDs []string

	for _, l := range logs {
		if l.ActorType == "user" && l.ActorID != nil {
			userIDs = append(userIDs, l.ActorID.String())
		}
		if (l.ActionCode == "ai.score_computed" || l.ActionCode == "ai.resume_parsed") && l.TargetID != nil {
			candidateIDs = append(candidateIDs, l.TargetID.String())
		}
	}

	usersMap := make(map[string]*pbAuth.UserObj)
	if len(userIDs) > 0 && u.authClient != nil && *u.authClient != nil {
		userRes, authErr := (*u.authClient).GetUsers(ctx, &pbAuth.GetUsersRequest{UserIds: userIDs})
		if authErr != nil {
			u.log.Error("failed to fetch users from auth service", zap.Error(authErr))
		} else {
			if userRes != nil && userRes.GetUsers() != nil {
				usersMap = userRes.GetUsers()
			}
		}
	}

	scoresMap := make(map[string]*pbAI.GetCandidateScoreResponse)
	if len(candidateIDs) > 0 && u.aiClient != nil && *u.aiClient != nil {
		for _, candidateID := range candidateIDs {
			scoreRes, aiErr := (*u.aiClient).GetCandidateScore(ctx, &pbAI.GetCandidateScoreRequest{CandidateId: candidateID})
			if aiErr != nil {
				u.log.Error("failed to fetch candidate scores from ai_engine", zap.Error(aiErr))
			} else {
				if scoreRes != nil && scoreRes.GetHasScore() {
					scoresMap[candidateID] = scoreRes
				}
			}
		}
	}

	for i := range logs {
		l := &logs[i]

		if l.ActorType == "user" && l.ActorID != nil {
			if u, ok := usersMap[l.ActorID.String()]; ok {
				name := u.FirstName
				if u.LastName != "" {
					name += " " + u.LastName
				}
				if name == "" {
					name = u.Email
				}
				l.ActorName = name
			} else {
				l.ActorName = "User"
			}
		} else if l.ActorType == "ai_agent" {
			l.ActorName = "AI Agent"
		} else {
			l.ActorName = "System"
		}

		if (l.ActionCode == "ai.score_computed" || l.ActionCode == "ai.resume_parsed") && l.TargetID != nil {
			if s, ok := scoresMap[l.TargetID.String()]; ok && s.HasScore {
				score := s.MatchScore
				l.MatchScore = &score
			}
		}
	}

	return logs, nil
}
