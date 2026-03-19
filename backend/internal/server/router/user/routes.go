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
	GetMembers() echo.HandlerFunc
	GetProfile() echo.HandlerFunc
	PatchProfile() echo.HandlerFunc
	PatchPassword() echo.HandlerFunc
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
		router.NewRoute(http.MethodGet, "/members", r.handler.GetMembers, r.rateLimit, r.session),
		router.NewRoute(http.MethodGet, "/profile", r.handler.GetProfile, r.rateLimit, r.session),
		router.NewRoute(http.MethodPatch, "/profile", r.handler.PatchProfile, r.rateLimit, r.session),
		router.NewRoute(http.MethodPatch, "/password", r.handler.PatchPassword, r.rateLimit, r.session),
	}
}
