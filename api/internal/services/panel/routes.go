package panel

import (
	"api/internal/routes/middlewares"
	"github.com/labstack/echo/v4"
)

func Routes(e *echo.Group) {
	controller := New()
	panelGroup := e.Group(
		"/panel",
		middlewares.IsAuthenticatedMiddleware(),
		middlewares.IsAdminPermissionMiddleware(),
	)
	panelGroup.PATCH("/config", controller.UpdatePanelConfig)
}
