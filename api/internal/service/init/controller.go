package init

import (
	"api/internal/repository"
	"api/pkg/utils"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"net/http"
)

type Controller struct {
	validator *validator.Validate
	adminRepo repository.AdminRepositoryInterface
}

func New() *Controller {
	return &Controller{
		validator: validator.New(),
		adminRepo: repository.NewAdminRepository(),
	}
}

func (ctrl *Controller) CheckToken(c echo.Context) error {
	var data struct {
		Token string `json:"token"`
	}
	if err := c.Bind(&data); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}
	if err := ctrl.validator.Struct(&data); err != nil {
		return c.JSON(http.StatusBadRequest, utils.InvalidBodyError(err))
	}

	//TODO: find best way to check token
	return c.JSON(http.StatusOK, nil)
}
