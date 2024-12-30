package routes

import (
	"api/internal/services/initialize"
	ocGroup "api/internal/services/oc_group"
	"api/internal/services/panel"
	"api/internal/services/staff"
	"api/internal/services/user"
	"github.com/labstack/echo/v4"
)

func Register(e *echo.Echo) {
	group := e.Group("/services/v1")
	initialize.Routes(group)
	panel.Routes(group)
	user.Routes(group)
	staff.Routes(group)
	ocGroup.Routes(group)
}
