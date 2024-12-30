package middlewares

import (
	"github.com/labstack/echo/v4"
)

func IsAdminPermissionMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		// TODO: check is admin
		return func(c echo.Context) error {
			// TODO: get userID and permission from context, then check user has this routes permission or not
			return next(c)
		}
	}
}
