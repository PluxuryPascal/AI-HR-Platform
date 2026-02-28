package invite

import (
	"backend/pkg/router"
	"net/http"

	"github.com/labstack/echo/v4"
)

type InviteRoutes interface {
	PostInvite() echo.HandlerFunc
	PostValidate() echo.HandlerFunc
	PostCreateUser() echo.HandlerFunc
}

type inviteRouter struct {
	routes     []router.Route
	handler    InviteRoutes
	middleware echo.MiddlewareFunc
}

func (r *inviteRouter) Routes() []router.Route {
	return r.routes
}

var _ router.Router = (*inviteRouter)(nil)

func NewRouter(h InviteRoutes, m echo.MiddlewareFunc) router.Router {
	r := &inviteRouter{
		handler:    h,
		middleware: m,
	}

	r.initRoutes()

	return r
}

func (r *inviteRouter) initRoutes() {
	r.routes = []router.Route{
		router.NewRoute(http.MethodPost, "/invite", r.handler.PostInvite, r.middleware),
		router.NewRoute(http.MethodPost, "/validate", r.handler.PostValidate, r.middleware),
		router.NewRoute(http.MethodPost, "/create-user", r.handler.PostCreateUser, r.middleware),
	}
}
