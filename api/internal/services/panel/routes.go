package panel

import (
	"api/internal/routes/middlewares"
	"github.com/labstack/echo/v4"
)

func Routes(e *echo.Group) {
	controller := New()
	e.GET("/panel/config/init", controller.GetPanelInitConfig)
	panelGroup := e.Group(
		"/panel",
		middlewares.IsAuthenticatedMiddleware(),
		middlewares.IsAdminPermissionMiddleware(),
	)
	panelGroup.POST("/config", controller.CreatePanelConfig)
	panelGroup.PATCH("/config", controller.UpdatePanelConfig)
	panelGroup.GET("/config", controller.GetPanelConfig)
}
