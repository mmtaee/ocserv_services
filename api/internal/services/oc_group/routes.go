package oc_group

import (
	"api/internal/routes/middlewares"
	"github.com/labstack/echo/v4"
)

func Routes(e *echo.Group) {
	controller := New()
	group := e.Group("/ocserv/group")
	group.POST("",
		controller.UpdateDefaultOcservGroup,
		middlewares.IsAuthenticatedMiddleware(),
		middlewares.IsAdminPermissionMiddleware(),
	)

}
