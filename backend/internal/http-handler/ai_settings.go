package handler

import (
	"backend/internal/domain"
	"backend/internal/openrouter"
	"backend/internal/repo"
	"backend/internal/response"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type AiSettingsHandler struct {
	log          *zap.Logger
	aiRepo       repo.AiSettingsRepository
	openRouter   *openrouter.Client
}

func NewAiSettingsHandler(
	log *zap.Logger,
	aiRepo repo.AiSettingsRepository,
	openRouter *openrouter.Client,
) *AiSettingsHandler {
	return &AiSettingsHandler{
		log:          log,
		aiRepo:       aiRepo,
		openRouter:   openRouter,
	}
}

type AiSettingsResponse struct {
	TeamID     string  `json:"team_id"`
	APIKey     *string `json:"api_key"`
	ParseModel *string `json:"parse_model"`
	ScoreModel *string `json:"score_model"`
	EmbedModel *string `json:"embed_model"`
	ChatModel  *string `json:"chat_model"`
}

type FieldUpdateRequest struct {
	Value string `json:"value"`
}

type ApiKeyUpdateRequest struct {
	ApiKey string `json:"api_key"`
}

func (h *AiSettingsHandler) GetSettings() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		teamID, ok := c.Get("team_id").(string)
		if !ok || teamID == "" {
			return response.Error(c, http.StatusUnauthorized, "unauthorized")
		}

		settings, err := h.aiRepo.GetByTeamID(ctx, teamID)
		if err != nil {
			h.log.Error("failed to get ai settings", zap.Error(err))
			return response.Error(c, http.StatusInternalServerError, "internal server error")
		}

		if settings == nil {
			return response.OK(c, AiSettingsResponse{
				TeamID: teamID,
			})
		}

		var maskedApiKey *string
		if settings.APIKey != nil && *settings.APIKey != "" {
			masked := "sk-or-v1-********"
			maskedApiKey = &masked
		}

		return response.OK(c, AiSettingsResponse{
			TeamID:     settings.TeamID,
			APIKey:     maskedApiKey,
			ParseModel: settings.ParseModel,
			ScoreModel: settings.ScoreModel,
			EmbedModel: settings.EmbedModel,
			ChatModel:  settings.ChatModel,
		})
	}
}

func (h *AiSettingsHandler) UpdateApiKey() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		teamID, ok := c.Get("team_id").(string)
		if !ok || teamID == "" {
			return response.Error(c, http.StatusUnauthorized, "unauthorized")
		}

		var req ApiKeyUpdateRequest
		if err := c.Bind(&req); err != nil {
			return response.Error(c, http.StatusBadRequest, "invalid param")
		}

		existing, err := h.aiRepo.GetByTeamID(ctx, teamID)
		if err != nil {
			h.log.Error("failed to fetch existing ai settings", zap.Error(err))
			return response.Error(c, http.StatusInternalServerError, "internal server error")
		}

		// Handle masked API keys.
		apiKeyToSave := &req.ApiKey
		if req.ApiKey == "sk-or-v1-********" {
			if existing == nil || existing.APIKey == nil {
				return response.Error(c, http.StatusBadRequest, "cannot use masked key for new setting")
			}
			apiKeyToSave = existing.APIKey
		}

		settings := &domain.TeamAISettings{
			TeamID: teamID,
			APIKey: apiKeyToSave,
		}
		if existing != nil {
			settings = existing
			settings.APIKey = apiKeyToSave
		}

		if err := h.aiRepo.Upsert(ctx, settings); err != nil {
			h.log.Error("failed to update ai settings", zap.Error(err))
			return response.Error(c, http.StatusInternalServerError, "internal server error")
		}

		return response.OK(c, map[string]string{"status": "success"})
	}
}

func (h *AiSettingsHandler) UpdateParseModel() echo.HandlerFunc {
	return h.updateStringField(func(s *domain.TeamAISettings, val string) {
		s.ParseModel = &val
	})
}

func (h *AiSettingsHandler) UpdateScoreModel() echo.HandlerFunc {
	return h.updateStringField(func(s *domain.TeamAISettings, val string) {
		s.ScoreModel = &val
	})
}

func (h *AiSettingsHandler) UpdateEmbedModel() echo.HandlerFunc {
	return h.updateStringField(func(s *domain.TeamAISettings, val string) {
		s.EmbedModel = &val
	})
}

func (h *AiSettingsHandler) UpdateChatModel() echo.HandlerFunc {
	return h.updateStringField(func(s *domain.TeamAISettings, val string) {
		s.ChatModel = &val
	})
}

func (h *AiSettingsHandler) updateStringField(updateFn func(*domain.TeamAISettings, string)) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		teamID, ok := c.Get("team_id").(string)
		if !ok || teamID == "" {
			return response.Error(c, http.StatusUnauthorized, "unauthorized")
		}

		var req FieldUpdateRequest
		if err := c.Bind(&req); err != nil {
			return response.Error(c, http.StatusBadRequest, "invalid param")
		}

		existing, err := h.aiRepo.GetByTeamID(ctx, teamID)
		if err != nil {
			h.log.Error("failed to fetch existing ai settings", zap.Error(err))
			return response.Error(c, http.StatusInternalServerError, "internal server error")
		}

		settings := &domain.TeamAISettings{TeamID: teamID}
		if existing != nil {
			settings = existing
		}

		updateFn(settings, req.Value)

		if err := h.aiRepo.Upsert(ctx, settings); err != nil {
			h.log.Error("failed to update ai settings", zap.Error(err))
			return response.Error(c, http.StatusInternalServerError, "internal server error")
		}

		return response.OK(c, map[string]string{"status": "success"})
	}
}

func (h *AiSettingsHandler) GetModels() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		models, err := h.openRouter.GetAvailableModels(ctx)
		if err != nil {
			h.log.Error("failed to get openrouter models", zap.Error(err))
			return response.Error(c, http.StatusInternalServerError, "internal server error")
		}

		return response.OK(c, models)
	}
}
