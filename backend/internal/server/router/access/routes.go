package access

import (
	"backend/pkg/router"
	"net/http"

	"github.com/labstack/echo/v4"
)

type AccessRoutes interface {
	PostGrantAccess() echo.HandlerFunc
	DeleteRevokeAccess() echo.HandlerFunc
	GetJobAccessList() echo.HandlerFunc
}

type accessRouter struct {
	handler AccessRoutes
	session echo.MiddlewareFunc
	rbac    echo.MiddlewareFunc
}

func NewRouter(h AccessRoutes, session, rbac echo.MiddlewareFunc) router.Router {
	return &accessRouter{
		handler: h,
		session: session,
		rbac:    rbac,
	}
}

func (r *accessRouter) Routes() []router.Route {
	return []router.Route{
		router.NewRoute(http.MethodGet, "", r.handler.GetJobAccessList, r.session, r.rbac),
		router.NewRoute(http.MethodPost, "", r.handler.PostGrantAccess, r.session, r.rbac),
		router.NewRoute(http.MethodDelete, "/:user_id", r.handler.DeleteRevokeAccess, r.session, r.rbac),
	}
}
