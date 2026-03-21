package job

import (
	"backend/pkg/router"
	"net/http"

	"github.com/labstack/echo/v4"
)

type JobRoutes interface {
	PostJob() echo.HandlerFunc
	PostJobList() echo.HandlerFunc
	GetJob() echo.HandlerFunc
	PatchJob() echo.HandlerFunc
	PostJobPublish() echo.HandlerFunc
	PostJobClose() echo.HandlerFunc
	PostJobArchive() echo.HandlerFunc
	DeleteJob() echo.HandlerFunc
}

type jobRouter struct {
	routes  []router.Route
	handler JobRoutes
	session echo.MiddlewareFunc
	rbac    echo.MiddlewareFunc
}

func (r *jobRouter) Routes() []router.Route {
	return r.routes
}

var _ router.Router = (*jobRouter)(nil)

func NewRouter(h JobRoutes, session, rbac echo.MiddlewareFunc) router.Router {
	r := &jobRouter{
		handler: h,
		session: session,
		rbac:    rbac,
	}

	r.initRoutes()

	return r
}

func (r *jobRouter) initRoutes() {
	r.routes = []router.Route{
		router.NewRoute(http.MethodPost, "", r.handler.PostJob, r.session, r.rbac),
		router.NewRoute(http.MethodPost, "/list", r.handler.PostJobList, r.session, r.rbac),
		router.NewRoute(http.MethodGet, "/:id", r.handler.GetJob, r.session, r.rbac),
		router.NewRoute(http.MethodPatch, "/:id", r.handler.PatchJob, r.session, r.rbac),
		router.NewRoute(http.MethodDelete, "/:id", r.handler.DeleteJob, r.session, r.rbac),
		router.NewRoute(http.MethodPost, "/:id/publish", r.handler.PostJobPublish, r.session, r.rbac),
		router.NewRoute(http.MethodPost, "/:id/close", r.handler.PostJobClose, r.session, r.rbac),
		router.NewRoute(http.MethodPost, "/:id/archive", r.handler.PostJobArchive, r.session, r.rbac),
	}
}
