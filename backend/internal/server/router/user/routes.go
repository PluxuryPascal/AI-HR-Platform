package user

import (
	"backend/internal/handler"
	"backend/pkg/router"
	"net/http"

	"github.com/labstack/echo/v4"
)

type userRouter struct {
	routes     []router.Route
	handler    *handler.AuthHandler
	middleware echo.MiddlewareFunc
}

func (r *userRouter) Routes() []router.Route {
	return r.routes
}

var _ router.Router = (*userRouter)(nil)

func NewRouter(h *handler.AuthHandler, m echo.MiddlewareFunc) router.Router {
	r := &userRouter{
		handler:    h,
		middleware: m,
	}

	r.initRoutes()

	return r
}

func (r *userRouter) initRoutes() {
	r.routes = []router.Route{
		router.NewRoute(http.MethodPost, "/login", r.handler.PostLogin, r.middleware),
		router.NewRoute(http.MethodPost, "/register", r.handler.PostRegister, r.middleware),
		router.NewRoute(http.MethodPost, "/logout", r.handler.PostLogout, r.middleware),
	}
}
