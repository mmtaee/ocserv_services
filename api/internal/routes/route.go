package routes

import (
	"api/internal/services/initialize"
	"github.com/labstack/echo/v4"
)

func Register(e *echo.Echo) {
	group := e.Group("/api/v1")
	initialize.Routes(group)
}
