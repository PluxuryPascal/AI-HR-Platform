package candidate

import (
	"backend/pkg/router"
	"net/http"

	"github.com/labstack/echo/v4"
)

type CandidateRoutes interface {
	PostCandidate() echo.HandlerFunc
	PostCandidateList() echo.HandlerFunc
	GetCandidate() echo.HandlerFunc
	PatchCandidate() echo.HandlerFunc
	PostCandidateMove() echo.HandlerFunc
	DeleteCandidate() echo.HandlerFunc
}

type candidateRouter struct {
	handler CandidateRoutes
	session echo.MiddlewareFunc
	rbac    echo.MiddlewareFunc
}

func NewRouter(h CandidateRoutes, session, rbac echo.MiddlewareFunc) router.Router {
	return &candidateRouter{
		handler: h,
		session: session,
		rbac:    rbac,
	}
}

func (r *candidateRouter) Routes() []router.Route {
	return []router.Route{
		router.NewRoute(http.MethodGet, "/:id", r.handler.GetCandidate, r.session, r.rbac),
		router.NewRoute(http.MethodPatch, "/:id", r.handler.PatchCandidate, r.session, r.rbac),
		router.NewRoute(http.MethodDelete, "/:id", r.handler.DeleteCandidate, r.session, r.rbac),
		router.NewRoute(http.MethodPost, "/:id/move", r.handler.PostCandidateMove, r.session, r.rbac),
	}
}

type candidateJobRouter struct {
	*candidateRouter
}

func NewJobScopedRouter(h CandidateRoutes, session, rbac echo.MiddlewareFunc) router.Router {
	return &candidateJobRouter{
		candidateRouter: &candidateRouter{
			handler: h,
			session: session,
			rbac:    rbac,
		},
	}
}

func (r *candidateJobRouter) Routes() []router.Route {
	return []router.Route{
		router.NewRoute(http.MethodPost, "", r.handler.PostCandidate, r.session, r.rbac),
		router.NewRoute(http.MethodPost, "/list", r.handler.PostCandidateList, r.session, r.rbac),
	}
}
