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
				log.Println(err)
				return echo.ErrForbidden
			}
			secretKey := c.QueryParam("init_secret")
			if secretKey == "" {
				log.Println("secret key is empty")
				return echo.ErrForbidden
			}
			content, err := os.ReadFile(config.GetApp().InitSecretFile)
			if err != nil {
				log.Println(err)
				return c.JSON(http.StatusNotFound, nil)
			}
			if secretKey != string(content) {
				log.Println("invalid secret key")
				return echo.ErrForbidden
			}
			return next(c)
		}
	}
}

func Routes(e *echo.Group) {
	controller := New()
	group := e.Group("/init")
	group.GET("/check", controller.CheckSecretKey)
	group.POST("/admin", controller.CreateSuperUser)
	group.POST("/config", controller.PanelConfig)

}
