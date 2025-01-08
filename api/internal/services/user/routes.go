package user

import (
	"api/internal/routes/middlewares"
	"github.com/labstack/echo/v4"
)

func Routes(e *echo.Group) {
	controller := New()
	group := e.Group("/user")
	group.POST("/admin", controller.CreateSuperUser)
	group.POST("/login", controller.Login)
	group.DELETE("/logout", controller.Logout, middlewares.IsAuthenticatedMiddleware())
	group.POST("/change_password", controller.ChangePassword, middlewares.IsAuthenticatedMiddleware())
}
