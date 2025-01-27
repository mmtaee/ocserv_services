package ocGroup

import (
	"api/internal/routes/middlewares"
	"github.com/labstack/echo/v4"
)

func Routes(e *echo.Group) {
	controller := New()
	group := e.Group("/ocserv/groups", middlewares.IsAuthenticatedMiddleware())

	group.POST("/defaults", controller.UpdateDefaultOcservGroup, middlewares.IsAdminPermissionMiddleware())
	group.GET("/defaults", controller.DefaultGroup)

	group.GET("", controller.Groups)
	group.POST("", controller.CreateGroup)
	group.GET("/:name", controller.Group)
	group.PATCH("/:name", controller.UpdateGroup)
	group.DELETE("/:name", controller.DeleteGroup)

	group.GET("/names", controller.GroupNames)
}
