package middlewares

import (
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/mmtaee/go-oc-utils/database"
	"github.com/mmtaee/go-oc-utils/models"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Unauthorized struct {
	Error string `json:"error"`
}

func unauthorized(c echo.Context) error {
	return c.JSON(http.StatusUnauthorized, Unauthorized{Error: "invalid authentication credentials"})
}

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
				return unauthorized(c)
			} else if err != nil {
				return c.JSON(http.StatusInternalServerError, nil)
			}

			c.Set("userID", strconv.Itoa(int(token.User.ID)))
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
