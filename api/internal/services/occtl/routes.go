package occtl

import (
	"api/internal/routes/middlewares"
	"github.com/labstack/echo/v4"
)

func Routes(e *echo.Group) {
	controller := New()

	group := e.Group("/occtl", middlewares.IsAuthenticatedMiddleware())

	group.POST("/reload", controller.Reload)
	group.GET("/online", controller.OnlineUsers)
	group.POST("/disconnect/:username", controller.Disconnect)
	group.GET("/ip_bans", controller.ShowIPBans)
	group.GET("/ip_bans/point", controller.ShowIPBansPoint)
	group.POST("/unban", controller.UnBanIP)
	group.GET("/status", controller.ShowStatus)
	//group.GET("/iroutes", controller.ShowIRoutes)
	group.GET("/users/:username", controller.ShowUser)
}
