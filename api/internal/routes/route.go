package routes

import (
	"api/internal/routes/middlewares"
	"api/internal/services/events"
	ocGroup "api/internal/services/oc_group"
	ocUser "api/internal/services/oc_user"
	"api/internal/services/occtl"
	"api/internal/services/panel"
	staffManagement "api/internal/services/staff_management"
	"api/internal/services/statistics"
	"api/internal/services/user"
	"github.com/labstack/echo/v4"
)

func Register(e *echo.Echo) {
	group := e.Group("/api/v1", middlewares.NeedInitMiddleware())
	panel.Routes(group)
	user.Routes(group)
	staffManagement.Routes(group)
	ocGroup.Routes(group)
	ocUser.Routes(group)
	statistics.Routes(group)
	occtl.Routes(group)
	events.Routes(group)
}
