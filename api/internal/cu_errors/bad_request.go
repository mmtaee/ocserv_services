package cu_errors

import (
	"errors"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/labstack/echo/v4"
	"net/http"
)

type ErrorResponse struct {
	Error string `json:"error" validate:"required"`
}

func BadRequest(c echo.Context, err error) error {
	var pqErr *pgconn.PgError
	switch {
	case errors.As(err, &pqErr):
		return c.JSON(http.StatusBadRequest, err.(*pgconn.PgError).Detail)
	default:
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}
}
