package workflow

import (
	"errors"
	"testing"

	"backend/internal/domain"
	hiringv1 "backend/internal/proto/hiring/v1"
	"backend/internal/temporal/activity"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/testsuite"
)

var pipelineInput = activity.ResumePipelineInput{
	CandidateID:   "cand-1",
	JobID:         "job-1",
	ResumeFileKey: "resumes/cand-1.pdf",
	TeamID:        "team-1",
	Locale:        "ru",
}

// newPipelineEnv поднимает тестовое окружение Temporal и мокает все активности,
// которых тест не касается явно.
func newPipelineEnv(t *testing.T) (*testsuite.TestWorkflowEnvironment, *activity.Activities) {
	t.Helper()

	var s testsuite.WorkflowTestSuite
	env := s.NewTestWorkflowEnvironment()

	var a *activity.Activities
	env.RegisterActivity(a)

	return env, a
}

func strPtr(s string) *string { return &s }

func TestResumePipeline_Success(t *testing.T) {
	env, a := newPipelineEnv(t)

	env.OnActivity(a.PDFExtract, mock.Anything, pipelineInput).Return(strPtr("Иван Иванов, Go-разработчик"), nil).Once()
	env.OnActivity(a.LLMParse, mock.Anything, mock.MatchedBy(func(in activity.LLMParseInput) bool {
		return in.CandidateID == "cand-1" && in.ParsedText == "Иван Иванов, Go-разработчик" && in.Locale == "ru"
	})).Return(&domain.ParseResult{JobRequirements: "Go, PostgreSQL"}, nil).Once()
	env.OnActivity(a.LLMScore, mock.Anything, mock.MatchedBy(func(in activity.LLMScoreInput) bool {
		return in.JobRequirements == "Go, PostgreSQL"
	})).Return(&domain.ScoreResult{MatchScore: 87}, nil).Once()
	env.OnActivity(a.LLMEmbed, mock.Anything, mock.Anything).Return(nil).Once()
	env.OnActivity(a.GRPCCallback, mock.Anything, mock.Anything, mock.Anything).Return(nil).Maybe()

	env.ExecuteWorkflow(ResumePipelineWorkflow, pipelineInput)

	require.True(t, env.IsWorkflowCompleted())
	require.NoError(t, env.GetWorkflowError())
	env.AssertExpectations(t)
}

func TestResumePipeline_NonRetryablePDFErrorIsNotRetried(t *testing.T) {
	env, a := newPipelineEnv(t)

	pdfErr := temporal.NewApplicationError("pdf is password protected", "NonRetryablePDFError")
	env.OnActivity(a.PDFExtract, mock.Anything, mock.Anything).Return(nil, pdfErr)
	env.OnActivity(a.GRPCCallback, mock.Anything, mock.Anything, mock.Anything).Return(nil).Maybe()

	env.ExecuteWorkflow(ResumePipelineWorkflow, pipelineInput)

	require.True(t, env.IsWorkflowCompleted())
	err := env.GetWorkflowError()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to extract pdf")

	// Ошибка помечена как неретраебельная — активность вызвана ровно один раз.
	env.AssertNumberOfCalls(t, "PDFExtract", 1)
	env.AssertNotCalled(t, "LLMParse", mock.Anything, mock.Anything)
}

func TestResumePipeline_ScoreFailureIsRetriedAndFailsWorkflow(t *testing.T) {
	env, a := newPipelineEnv(t)

	env.OnActivity(a.PDFExtract, mock.Anything, mock.Anything).Return(strPtr("text"), nil)
	env.OnActivity(a.LLMParse, mock.Anything, mock.Anything).Return(&domain.ParseResult{}, nil)
	env.OnActivity(a.LLMScore, mock.Anything, mock.Anything).Return(nil, errors.New("openrouter: 503"))
	env.OnActivity(a.LLMEmbed, mock.Anything, mock.Anything).Return(nil)
	env.OnActivity(a.GRPCCallback, mock.Anything, mock.Anything, mock.Anything).Return(nil).Maybe()

	env.ExecuteWorkflow(ResumePipelineWorkflow, pipelineInput)

	require.True(t, env.IsWorkflowCompleted())
	err := env.GetWorkflowError()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to score")

	// RetryPolicy для LLM: MaximumAttempts = 3.
	env.AssertNumberOfCalls(t, "LLMScore", 3)
	// Эмбеддинги считаются параллельно и не зависят от падения скоринга.
	env.AssertNumberOfCalls(t, "LLMEmbed", 1)
}

func TestDomainStatusToProto(t *testing.T) {
	tests := []struct {
		in   domain.CandidateParsingStatus
		want hiringv1.ParsingStatus
	}{
		{domain.ParsingStatusNeedsReview, hiringv1.ParsingStatus_PARSING_STATUS_NEEDS_REVIEW},
		{domain.ParsingStatusCompleted, hiringv1.ParsingStatus_PARSING_STATUS_SUCCESS},
		{domain.ParsingStatusPending, hiringv1.ParsingStatus_PARSING_STATUS_SUCCESS},
	}

	for _, tt := range tests {
		t.Run(string(tt.in), func(t *testing.T) {
			assert.Equal(t, tt.want, domainStatusToProto(tt.in))
		})
	}
}
