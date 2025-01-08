package routes

import (
	"api/internal/models"
	ocGroup "api/internal/services/oc_group"
	ocUser "api/internal/services/oc_user"
	"api/internal/services/panel"
	staffManagement "api/internal/services/staff_management"
	"api/internal/services/statistics"
	"api/internal/services/user"
	"api/pkg/config"
	"api/pkg/database"
	"github.com/labstack/echo/v4"
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

func Register(e *echo.Echo) {
	group := e.Group("/api/v1", NeedInitMiddleware())
	panel.Routes(group)
	user.Routes(group)
	staffManagement.Routes(group)
	ocGroup.Routes(group)
	ocUser.Routes(group)
	statistics.Routes(group)
}
