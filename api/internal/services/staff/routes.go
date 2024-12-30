package staff

import (
	"api/internal/routes/middlewares"
	"github.com/labstack/echo/v4"
)

func Routes(e *echo.Group) {
	controller := New()
	staffGroup := e.Group(
		"/staff",
		middlewares.IsAuthenticatedMiddleware(),
		middlewares.IsAdminPermissionMiddleware(),
	)
	staffGroup.GET("", controller.Staffs)
	staffGroup.POST("", controller.CreateStaff)
	staffGroup.GET("/:id", controller.Staff)
	staffGroup.PATCH("/:id", controller.UpdateStaff)
	staffGroup.DELETE("/:id", controller.DeleteStaff)
}
