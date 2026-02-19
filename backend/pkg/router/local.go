package router

import (
	"context"

	"github.com/labstack/echo/v4"
)

type route struct {
	method      string
	path        string
	handler     func() echo.HandlerFunc
	middlewares []echo.MiddlewareFunc
}

func (r *route) Register(ctx context.Context, e *echo.Echo) {
	e.Add(r.method, r.path, r.handler(), r.middlewares...)
}

var _ Route = (*route)(nil)

func NewRoute(method, path string, handler func() echo.HandlerFunc, middlewares ...echo.MiddlewareFunc) Route {
	return &route{
		method:      method,
		path:        path,
		handler:     handler,
		middlewares: middlewares,
	}
}
