package response

import (
	"math"
	"net/http"

	"github.com/labstack/echo/v4"
)

type ApiResponse[T any] struct {
	Data    T      `json:"data"`
	Message string `json:"message,omitempty"`
	Success bool   `json:"success"`
}

type PaginationMeta struct {
	Page       int `json:"page"`
	PerPage    int `json:"per_page"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

type PaginatedResponse[T any] struct {
	Data    []T            `json:"data"`
	Meta    PaginationMeta `json:"meta"`
	Success bool           `json:"success"`
}

func OK[T any](c echo.Context, data T) error {
	return c.JSON(http.StatusOK, ApiResponse[T]{
		Data:    data,
		Success: true,
	})
}

func Created[T any](c echo.Context, data T) error {
	return c.JSON(http.StatusCreated, ApiResponse[T]{
		Data:    data,
		Success: true,
	})
}

func Paginated[T any](c echo.Context, data []T, total, page, perPage int) error {
	totalPages := int(math.Ceil(float64(total) / float64(perPage)))
	if perPage <= 0 {
		totalPages = 0
	}

	return c.JSON(http.StatusOK, PaginatedResponse[T]{
		Data: data,
		Meta: PaginationMeta{
			Page:       page,
			PerPage:    perPage,
			Total:      total,
			TotalPages: totalPages,
		},
		Success: true,
	})
}

func NoContent(c echo.Context, status int) error {
	return c.NoContent(status)
}

func Error(c echo.Context, status int, message string) error {
	return c.JSON(status, ApiResponse[any]{
		Success: false,
		Message: message,
	})
}
