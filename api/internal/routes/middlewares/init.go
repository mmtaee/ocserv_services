package middlewares

import (
	"api/pkg/config"
	"github.com/labstack/echo/v4"
	"github.com/mmtaee/go-oc-utils/database"
	"github.com/mmtaee/go-oc-utils/models"
	"net/http"
	"slices"
)

func NeedInitMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if c.Path() == "/api/v1/panel/config/init" {
				return next(c)
			}
			if slices.Contains([]string{
				"/api/v1/user/admin",
				"/api/v1/panel/config",
			}, c.Path()) && c.Request().Method == "POST" {
				return next(c)
			}
			if config.GetAppInit() {
				return next(c)
			}
			p := models.PanelConfig{}
			db := database.Connection()
			err := db.WithContext(c.Request().Context()).First(&p).Error
			if err != nil || !p.Init {
				return c.JSON(http.StatusForbidden, map[string]string{"error": "initializing required"})
			}
			if p.Init {
				config.ActiveAppInit()
			}
			return next(c)
		}
	}
}
