package initialize

import (
	"api/internal/routes/middlewares"
	"github.com/labstack/echo/v4"
)

func Routes(e *echo.Group) {
	controller := New()
	group := e.Group("/init")
	group.POST("/admin", controller.CreateSuperUser)
	group.POST("/config", controller.PanelConfig,
		middlewares.IsAuthenticatedMiddleware(),
		middlewares.IsAdminPermissionMiddleware(),
	)
}
