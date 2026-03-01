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
	routes    []router.Route
	handler   UserRoutes
	rateLimit echo.MiddlewareFunc
	session   echo.MiddlewareFunc
}

func (r *userRouter) Routes() []router.Route {
	return r.routes
}

var _ router.Router = (*userRouter)(nil)

func NewRouter(h UserRoutes, rateLimit echo.MiddlewareFunc, session echo.MiddlewareFunc) router.Router {
	r := &userRouter{
		handler:   h,
		rateLimit: rateLimit,
		session:   session,
	}

	r.initRoutes()

	return r
}

func (r *userRouter) initRoutes() {
	r.routes = []router.Route{
		router.NewRoute(http.MethodPost, "/login", r.handler.PostLogin, r.rateLimit),
		router.NewRoute(http.MethodPost, "/register", r.handler.PostRegister, r.rateLimit),
		router.NewRoute(http.MethodPost, "/logout", r.handler.PostLogout, r.rateLimit, r.session),
	}
}
