package outreach

import (
	"backend/pkg/router"
	"net/http"

	"github.com/labstack/echo/v4"
)

type OutreachRoutes interface {
	PostGenerateEmail() echo.HandlerFunc
	PostSendEmail() echo.HandlerFunc
}

type outreachRouter struct {
	handler OutreachRoutes
	session echo.MiddlewareFunc
	rbac    echo.MiddlewareFunc
}

func NewRouter(h OutreachRoutes, session, rbac echo.MiddlewareFunc) router.Router {
	return &outreachRouter{
		handler: h,
		session: session,
		rbac:    rbac,
	}
}

func (r *outreachRouter) Routes() []router.Route {
	return []router.Route{
		router.NewRoute(http.MethodPost, "/:candidate_id/generate", r.handler.PostGenerateEmail, r.session, r.rbac),
		router.NewRoute(http.MethodPost, "/:id/send", r.handler.PostSendEmail, r.session, r.rbac),
	}
}
