package utils

import (
	"errors"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/labstack/echo/v4"
	"net/http"
)

type ErrorResponse struct {
	Error string `json:"error" validate:"required"`
}

func BadRequest(c echo.Context, err interface{}) error {
	switch err.(type) {
	case error:
		var pqErr *pgconn.PgError
		if errors.As(err.(error), &pqErr) {
			return c.JSON(http.StatusBadRequest, err.(*pgconn.PgError).Detail)
		}
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.(error).Error(),
		})
	case string:
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.(string),
		})
	case map[string]interface{}:
		return c.JSON(http.StatusBadRequest, errors.New(err.(map[string]interface{})["error"].([]string)[0]))
	case interface{}:
		return c.JSON(http.StatusBadRequest, err.(interface{}))
	default:
		return c.JSON(http.StatusBadRequest, errors.New(http.StatusText(http.StatusBadRequest)))
	}
}
