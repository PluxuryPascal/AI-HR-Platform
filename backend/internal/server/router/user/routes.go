package user

import (
	"backend/pkg/router"
	"net/http"

	"github.com/labstack/echo/v4"
)

type UserRoutes interface {
	PostLogin() echo.HandlerFunc
	PostRegister() echo.HandlerFunc
	PostLogout() echo.HandlerFunc
}

type userRouter struct {
	routes     []router.Route
	handler    UserRoutes
	middleware echo.MiddlewareFunc
}

func (r *userRouter) Routes() []router.Route {
	return r.routes
}

var _ router.Router = (*userRouter)(nil)

func NewRouter(h UserRoutes, m echo.MiddlewareFunc) router.Router {
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
