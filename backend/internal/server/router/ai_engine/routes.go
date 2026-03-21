package ai_engine

import (
	"backend/pkg/router"
	"net/http"

	"github.com/labstack/echo/v4"
)

type AiRoutes interface {
	ParseJob() echo.HandlerFunc
}

type aiRouter struct {
	handler AiRoutes
	session echo.MiddlewareFunc
	rbac    echo.MiddlewareFunc
}

func NewRouter(h AiRoutes, session, rbac echo.MiddlewareFunc) router.Router {
	return &aiRouter{
		handler: h,
		session: session,
		rbac:    rbac,
	}
}

func (r *aiRouter) Routes() []router.Route {
	return []router.Route{
		router.NewRoute(http.MethodPost, "/parse-job", r.handler.ParseJob, r.session, r.rbac),
	}
}
