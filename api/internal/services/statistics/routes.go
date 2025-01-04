package statistics

import "github.com/labstack/echo/v4"

func Routes(e *echo.Group) {
	controller := New()
	group := e.Group("/statistics")
	group.GET("/", controller.Statistics)
}
