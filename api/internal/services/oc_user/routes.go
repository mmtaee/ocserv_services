package ocUser

import "github.com/labstack/echo/v4"

func Routes(e *echo.Group) {
	controller := New()
	group := e.Group("/ocserv/users")
	group.GET("", controller.Users)
	group.POST("", controller.Create)
	group.GET("/:id", controller.User)
	group.PATCH("/:id", controller.Update)
	group.POST("/:id/lock", controller.Lock)
	group.POST("/:id/disconnect", controller.Disconnect)
	group.DELETE("/:id", controller.Delete)
	group.GET("/:id/activity", controller.Activity)
	group.GET("/:id/statistics", controller.Statistics)
}
