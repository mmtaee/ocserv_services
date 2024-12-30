package initialize

import (
	"github.com/labstack/echo/v4"
)

func Routes(e *echo.Group) {
	controller := New()
	group := e.Group("/init")
	group.GET("/check", controller.CheckSecretKey)
	group.POST("/admin", controller.CreateSuperUser)
	group.POST("/config", controller.PanelConfig)
}
