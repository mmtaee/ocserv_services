package routes

import (
	"api/internal/services/initialize"
	"api/internal/services/user"
	"github.com/labstack/echo/v4"
)

func Register(e *echo.Echo) {
	group := e.Group("/api/v1")
	initialize.Routes(group)
	user.Routes(group)
}
