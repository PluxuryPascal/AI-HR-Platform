package activity

import (
	"context"
	"fmt"

	hiringv1 "backend/internal/proto/hiring/v1"

	"go.uber.org/zap"
)

func (a *Activities) InterviewGetCandidateText(ctx context.Context, input InterviewInput) (string, error) {
	a.log.Debug("InterviewGetCandidateText started", zap.String("candidate_id", input.CandidateID))

	conn := a.hiringGRPC.GetConn()
	client := hiringv1.NewHiringServiceClient(conn)

	resp, err := client.GetCandidate(ctx, &hiringv1.GetCandidateRequest{
		CandidateId: input.CandidateID,
	})
	if err != nil {
		return "", fmt.Errorf("hiring grpc GetCandidate: %w", err)
	}

	return resp.GetParsedText(), nil
}

func (a *Activities) InterviewGenerateQuestions(ctx context.Context, input InterviewInput, resumeText string) (InterviewOutput, error) {
	a.log.Debug("InterviewGenerateQuestions started", zap.String("candidate_id", input.CandidateID))

	res, err := a.interviewGen.Generate(ctx, input.TeamID, resumeText, input.Locale)
	if err != nil {
		return InterviewOutput{}, fmt.Errorf("llm generate interview: %w", err)
	}

	questions := make([]InterviewPair, 0, len(res))
	for _, q := range res {
		questions = append(questions, InterviewPair{
			ID:       q.ID,
			Question: q.Question,
			Answer:   q.Answer,
			Category: q.Category,
		})
	}

	return InterviewOutput{Questions: questions}, nil
}
