package initialize

import (
	"api/pkg/config"
	"errors"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"os"
)

func IntiRoutePermissionMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			file := config.GetApp().InitSecretFile
			_, err := os.Stat(file)
			if err != nil {
				if !errors.Is(err, os.ErrNotExist) {
					log.Println(err)
				}
				return echo.ErrForbidden
			}
			secretKey := c.QueryParam("secret_key")
			if secretKey == "" {
				return echo.ErrForbidden
			}
			content, err := os.ReadFile(config.GetApp().InitSecretFile)
			if err != nil {
				return c.JSON(http.StatusNotFound, nil)
			}
			if secretKey != string(content) {
				return echo.ErrForbidden
			}
			return next(c)
		}
	}
}

func Routes(e *echo.Group) {
	controller := New()
	group := e.Group("/init", IntiRoutePermissionMiddleware())
	group.POST("/admin", controller.CreateSuperUser)
	group.POST("/config", controller.PanelConfig)
	group.POST("/group", controller.DefaultOcservGroup)
}
