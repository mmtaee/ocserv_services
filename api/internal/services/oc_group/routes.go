package ocGroup

import (
	"api/internal/routes/middlewares"
	"github.com/labstack/echo/v4"
)

func Routes(e *echo.Group) {
	controller := New()
	group := e.Group("/ocserv/groups")
	group.POST("",
		controller.UpdateDefaultOcservGroup,
		middlewares.IsAuthenticatedMiddleware(),
		middlewares.IsAdminPermissionMiddleware(),
	)

	group.GET("", controller.Groups)
	group.PATCH("/:name", controller.UpdateGroup)
	group.DELETE("/:name", controller.DeleteGroup)
}
