package dashboard

import (
	"backend/pkg/router"
	"net/http"

	"github.com/labstack/echo/v4"
)

type DashboardRoutes interface {
	GetStats() echo.HandlerFunc
	GetApplicationDynamics() echo.HandlerFunc
	GetRecentActivity() echo.HandlerFunc
}

type dashboardRouter struct {
	handler DashboardRoutes
	session echo.MiddlewareFunc
	rbac    echo.MiddlewareFunc
}

func NewRouter(h DashboardRoutes, session, rbac echo.MiddlewareFunc) router.Router {
	return &dashboardRouter{
		handler: h,
		session: session,
		rbac:    rbac,
	}
}

func (r *dashboardRouter) Routes() []router.Route {
	return []router.Route{
		router.NewRoute(http.MethodGet, "/stats", r.handler.GetStats, r.session),
		router.NewRoute(http.MethodGet, "/applications-chart", r.handler.GetApplicationDynamics, r.session),
		router.NewRoute(http.MethodGet, "/recent-activity", r.handler.GetRecentActivity, r.session),
	}
}
