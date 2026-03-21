package ai_settings

import (
	"backend/pkg/router"
	"net/http"

	"github.com/labstack/echo/v4"
)

type AiSettingsRoutes interface {
	GetSettings() echo.HandlerFunc
	UpdateApiKey() echo.HandlerFunc
	UpdateParseModel() echo.HandlerFunc
	UpdateScoreModel() echo.HandlerFunc
	UpdateEmbedModel() echo.HandlerFunc
	UpdateChatModel() echo.HandlerFunc
	GetModels() echo.HandlerFunc
}

type aiSettingsRouter struct {
	handler AiSettingsRoutes
	session echo.MiddlewareFunc
	rbac    echo.MiddlewareFunc
}

func NewRouter(h AiSettingsRoutes, session, rbac echo.MiddlewareFunc) router.Router {
	return &aiSettingsRouter{
		handler: h,
		session: session,
		rbac:    rbac,
	}
}

func (r *aiSettingsRouter) Routes() []router.Route {
	return []router.Route{
		router.NewRoute(http.MethodGet, "/models", r.handler.GetModels, r.session, r.rbac),
		router.NewRoute(http.MethodGet, "/", r.handler.GetSettings, r.session, r.rbac),
		router.NewRoute(http.MethodPatch, "/api-key", r.handler.UpdateApiKey, r.session, r.rbac),
		router.NewRoute(http.MethodPatch, "/parse-model", r.handler.UpdateParseModel, r.session, r.rbac),
		router.NewRoute(http.MethodPatch, "/score-model", r.handler.UpdateScoreModel, r.session, r.rbac),
		router.NewRoute(http.MethodPatch, "/embed-model", r.handler.UpdateEmbedModel, r.session, r.rbac),
		router.NewRoute(http.MethodPatch, "/chat-model", r.handler.UpdateChatModel, r.session, r.rbac),
	}
}
