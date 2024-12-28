package user

import (
	"api/internal/routes/middlewares"
	"github.com/labstack/echo/v4"
)

func Routes(e *echo.Group) {
	controller := New()
	group := e.Group("/user")
	group.POST("/login", controller.Login)
	group.DELETE("/logout", controller.Logout)
	group.POST("/password", controller.ChangePassword, middlewares.IsAuthenticatedMiddleware())

	staffGroup := e.Group(
		"/staff",
		middlewares.IsAuthenticatedMiddleware(),
		middlewares.IsAdminPermissionMiddleware(),
	)
	staffGroup.GET("/", controller.Staffs)
	staffGroup.POST("/", controller.CreateStaff)
	staffGroup.GET("/:id", controller.Staff)
	staffGroup.PATCH("/:id", controller.UpdateStaff)
	staffGroup.DELETE("/:id", controller.DeleteStaff)

	panelGroup := e.Group(
		"/panel",
		middlewares.IsAuthenticatedMiddleware(),
		middlewares.IsAdminPermissionMiddleware(),
	)
	panelGroup.PATCH("/config", controller.UpdatePanelConfig)
}
