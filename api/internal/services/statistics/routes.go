package statistics

import (
	"api/internal/routes/middlewares"
	"github.com/labstack/echo/v4"
)

func Routes(e *echo.Group) {
	controller := New()
	group := e.Group("/statistics", middlewares.IsAuthenticatedMiddleware())
	group.GET("", controller.Statistics)
}
