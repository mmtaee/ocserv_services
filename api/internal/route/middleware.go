package route

import (
	"github.com/labstack/echo/v4"
	"time"
)

func HasPermissionMiddleware(timeout time.Duration) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// TODO: get userID and permission from context, then check user has this route permission or not
			return next(c)
		}
	}
}
