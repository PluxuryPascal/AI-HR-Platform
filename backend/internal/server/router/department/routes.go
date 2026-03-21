package department

import (
	"backend/pkg/router"
	"net/http"

	"github.com/labstack/echo/v4"
)

type DepartmentRoutes interface {
	PostDepartment() echo.HandlerFunc
	GetDepartments() echo.HandlerFunc
	DeleteDepartment() echo.HandlerFunc
}

type departmentRouter struct {
	handler DepartmentRoutes
	session echo.MiddlewareFunc
	rbac    echo.MiddlewareFunc
}

func NewRouter(h DepartmentRoutes, session, rbac echo.MiddlewareFunc) router.Router {
	return &departmentRouter{
		handler: h,
		session: session,
		rbac:    rbac,
	}
}

func (r *departmentRouter) Routes() []router.Route {
	return []router.Route{
		router.NewRoute(http.MethodPost, "", r.handler.PostDepartment, r.session, r.rbac),
		router.NewRoute(http.MethodGet, "", r.handler.GetDepartments, r.session),
		router.NewRoute(http.MethodDelete, "/:id", r.handler.DeleteDepartment, r.session, r.rbac),
	}
}
