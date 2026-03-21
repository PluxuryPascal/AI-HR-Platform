package chat

import (
	"backend/pkg/router"
	"net/http"

	"github.com/labstack/echo/v4"
)

type ChatRoutes interface {
	PostChat() echo.HandlerFunc
	GetChatSessions() echo.HandlerFunc
	GetChatHistory() echo.HandlerFunc
}

type chatRouter struct {
	handler ChatRoutes
	session echo.MiddlewareFunc
	rbac    echo.MiddlewareFunc
}

func NewRouter(h ChatRoutes, session, rbac echo.MiddlewareFunc) router.Router {
	return &chatRouter{
		handler: h,
		session: session,
		rbac:    rbac,
	}
}

func (r *chatRouter) Routes() []router.Route {
	return []router.Route{
		router.NewRoute(http.MethodPost, "", r.handler.PostChat, r.session, r.rbac),
		router.NewRoute(http.MethodGet, "/sessions", r.handler.GetChatSessions, r.session),
		router.NewRoute(http.MethodGet, "/sessions/:id/history", r.handler.GetChatHistory, r.session),
	}
}
