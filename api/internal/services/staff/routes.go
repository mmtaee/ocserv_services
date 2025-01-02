package staff

import (
	"api/internal/routes/middlewares"
	"github.com/labstack/echo/v4"
)

func Routes(e *echo.Group) {
	controller := New()
	staffGroup := e.Group(
		"/staffs",
		middlewares.IsAuthenticatedMiddleware(),
		middlewares.IsAdminPermissionMiddleware(),
	)
	staffGroup.GET("", controller.Staffs)
	staffGroup.POST("", controller.CreateStaff)
	staffGroup.GET("/:uid/permission", controller.StaffPermission)
	staffGroup.PATCH("/:uid/permission", controller.UpdateStaffPermission)
	staffGroup.PATCH("/:uid", controller.UpdateStaffPermission)
	staffGroup.DELETE("/:uid", controller.DeleteStaff)
}
