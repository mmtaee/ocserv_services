package routes

import (
	"api/internal/services/initialize"
	ocGroup "api/internal/services/oc_group"
	ocUser "api/internal/services/oc_user"
	"api/internal/services/panel"
	staffManagement "api/internal/services/staff_management"
	"api/internal/services/statistics"
	"api/internal/services/user"
	"github.com/labstack/echo/v4"
)

func Register(e *echo.Echo) {
	group := e.Group("/api/v1")
	initialize.Routes(group)
	panel.Routes(group)
	user.Routes(group)
	staffManagement.Routes(group)
	ocGroup.Routes(group)
	ocUser.Routes(group)
	statistics.Routes(group)
}
