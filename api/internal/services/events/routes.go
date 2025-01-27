package events

import (
	"api/internal/routes/middlewares"
	"github.com/labstack/echo/v4"
)

func Routes(e *echo.Group) {
	controller := New()
	e.GET("/events/:event_type", controller.Events, middlewares.IsAuthenticatedMiddleware())
}
