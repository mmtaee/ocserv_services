package middlewares

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
)

type PermissionDenied struct {
	Error string `json:"error" validate:"required"`
}

func PermissionDeniedResponse(c echo.Context, msg ...string) error {
	if len(msg) == 0 {
		return c.JSON(http.StatusForbidden, PermissionDenied{
			Error: "Permission denied",
		})
	}
	return c.JSON(http.StatusForbidden, PermissionDenied{
		Error: strings.Join(msg, ", "),
	})
}

func IsAdminPermissionMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if isAdmin := c.Get("isAdmin"); !isAdmin.(bool) {
				return PermissionDeniedResponse(c, "admin permission required")
			}
			return next(c)
		}
	}
}
