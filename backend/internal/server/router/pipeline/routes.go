package pipeline

import (
	"backend/pkg/router"
	"net/http"

	"github.com/labstack/echo/v4"
)

type PipelineRoutes interface {
	PostStage() echo.HandlerFunc
	GetStages() echo.HandlerFunc
	PatchStage() echo.HandlerFunc
	DeleteStage() echo.HandlerFunc
}

type pipelineRouter struct {
	handler PipelineRoutes
	session echo.MiddlewareFunc
	rbac    echo.MiddlewareFunc
}

func NewRouter(h PipelineRoutes, session, rbac echo.MiddlewareFunc) router.Router {
	return &pipelineRouter{
		handler: h,
		session: session,
		rbac:    rbac,
	}
}

func (r *pipelineRouter) Routes() []router.Route {
	return []router.Route{
		router.NewRoute(http.MethodGet, "", r.handler.GetStages, r.session, r.rbac),
		router.NewRoute(http.MethodPost, "", r.handler.PostStage, r.session, r.rbac),
		router.NewRoute(http.MethodPatch, "/:stage_id", r.handler.PatchStage, r.session, r.rbac),
		router.NewRoute(http.MethodDelete, "/:stage_id", r.handler.DeleteStage, r.session, r.rbac),
	}
}
