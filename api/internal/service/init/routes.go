package init

import "github.com/labstack/echo/v4"

func Routes(e *echo.Echo) {
	controller := New()
	// TODO: add middleware to check has allow to init middleware
	group := e.Group("/init")
	group.POST("/token", controller.CheckToken)
}
