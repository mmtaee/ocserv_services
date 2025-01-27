package ocUser

import (
	"api/internal/routes/middlewares"
	"github.com/labstack/echo/v4"
)

func Routes(e *echo.Group) {
	controller := New()
	group := e.Group("/ocserv/users", middlewares.IsAuthenticatedMiddleware())

	group.GET("", controller.Users)
	group.POST("", controller.Create)
	group.GET("/:uid", controller.User)
	group.PATCH("/:uid", controller.Update)
	group.POST("/:uid/lock", controller.LockOrUnlock)
	group.POST("/:uid/disconnect", controller.Disconnect)
	group.DELETE("/:uid", controller.Delete)
	group.GET("/:uid/statistics", controller.Statistics)
	group.GET("/:uid/activities", controller.Activities)
}
