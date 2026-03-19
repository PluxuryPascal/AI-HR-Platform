package interview

import (
	"backend/pkg/router"
	"net/http"

	"github.com/labstack/echo/v4"
)

type InterviewRoutes interface {
	GenerateQuestions() echo.HandlerFunc
}

type interviewRouter struct {
	routes  []router.Route
	handler InterviewRoutes
	session echo.MiddlewareFunc
	rbac    echo.MiddlewareFunc
}

func (r *interviewRouter) Routes() []router.Route {
	return r.routes
}

var _ router.Router = (*interviewRouter)(nil)

func NewRouter(h InterviewRoutes, session, rbac echo.MiddlewareFunc) router.Router {
	r := &interviewRouter{
		handler: h,
		session: session,
		rbac:    rbac,
	}

	r.initRoutes()

	return r
}

func (r *interviewRouter) initRoutes() {
	r.routes = []router.Route{
		router.NewRoute(http.MethodPost, "/:id/questions", r.handler.GenerateQuestions, r.session, r.rbac),
	}
}
