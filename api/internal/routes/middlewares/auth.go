package middlewares

import (
	"api/internal/models"
	"api/pkg/database"
	"errors"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"net/http"
	"strings"
	"time"
)

func IsAuthenticatedMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token := models.UserToken{}
			tokenString := c.Request().Header.Get("Authorization")
			tokenString = strings.TrimPrefix(tokenString, "Bearer ")
			db := database.Connection()
			err := db.WithContext(c.Request().Context()).
				Table("user_tokens").
				Preload("User").Preload("User.Permission").
				Where("token = ? AND expire_at > ?", tokenString, time.Now()).
				First(&token).Error
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "invalid authentication credentials",
				})
			} else if err != nil {
				return c.JSON(http.StatusInternalServerError, nil)
			}
			c.Set("userID", token.User.ID)
			c.Set("username", token.User.Username)
			c.Set("isAdmin", token.User.IsAdmin)
			c.Set("token", token.Token)
			if !token.User.IsAdmin {
				c.Set("permission", token.User.Permission)
			}
			return next(c)
		}
	}
}
