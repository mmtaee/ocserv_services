package routes

import (
	"api/internal/services/initialize"
	ocGroup "api/internal/services/oc_group"
	"api/internal/services/user"
	"github.com/labstack/echo/v4"
)

func Register(e *echo.Echo) {
	group := e.Group("/api/v1")
	initialize.Routes(group)
	user.Routes(group)
	ocGroup.Routes(group)
}
