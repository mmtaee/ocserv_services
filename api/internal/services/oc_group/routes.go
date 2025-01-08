package ocGroup

import (
	"api/internal/routes/middlewares"
	"github.com/labstack/echo/v4"
)

func Routes(e *echo.Group) {
	controller := New()
	group := e.Group("/ocserv/groups")

	group.POST("/defaults",
		controller.UpdateDefaultOcservGroup,
		middlewares.IsAuthenticatedMiddleware(),
		middlewares.IsAdminPermissionMiddleware(),
	)

	group.GET("", controller.Groups)
	group.POST("", controller.CreateGroup)
	group.POST("/:name", controller.UpdateGroup)
	group.DELETE("/:name", controller.DeleteGroup)
}
