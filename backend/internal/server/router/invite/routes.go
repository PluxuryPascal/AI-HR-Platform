package invite

import (
	"backend/pkg/router"
	"net/http"

	"github.com/labstack/echo/v4"
)

type InviteRoutes interface {
	PostInvite() echo.HandlerFunc
	GetValidate() echo.HandlerFunc
	PostCreateUser() echo.HandlerFunc
}

type inviteRouter struct {
	routes    []router.Route
	handler   InviteRoutes
	rateLimit echo.MiddlewareFunc
	session   echo.MiddlewareFunc
	rbac      echo.MiddlewareFunc
}

func (r *inviteRouter) Routes() []router.Route {
	return r.routes
}

var _ router.Router = (*inviteRouter)(nil)

func NewRouter(h InviteRoutes, rateLimit echo.MiddlewareFunc, session echo.MiddlewareFunc, rbac echo.MiddlewareFunc) router.Router {
	r := &inviteRouter{
		handler:   h,
		rateLimit: rateLimit,
		session:   session,
		rbac:      rbac,
	}

	r.initRoutes()

	return r
}

func (r *inviteRouter) initRoutes() {
	r.routes = []router.Route{
		router.NewRoute(http.MethodPost, "/invite", r.handler.PostInvite, r.rateLimit, r.session, r.rbac),
		router.NewRoute(http.MethodGet, "/validate", r.handler.GetValidate, r.rateLimit),
		router.NewRoute(http.MethodPost, "/create-user", r.handler.PostCreateUser, r.rateLimit),
	}
}
