package ocUser

import "github.com/labstack/echo/v4"

func Routes(e *echo.Group) {
	controller := New()
	group := e.Group("/ocserv/users")
	group.GET("", controller.Users)
	group.POST("", controller.Create)
	group.GET("/:uid", controller.User)
	group.PATCH("/:uid", controller.Update)
	group.POST("/:uid/lock", controller.Lock)
	group.POST("/:uid/disconnect", controller.Disconnect)
	group.DELETE("/:uid", controller.Delete)
	group.GET("/:uid/activity", controller.Activity)
	group.GET("/:uid/statistics", controller.Statistics)
}
