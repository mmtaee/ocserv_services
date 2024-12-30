package panel

import (
	"api/pkg/validator"
	"github.com/labstack/echo/v4"
)

type Controller struct {
	validator validator.CustomValidatorInterface
}

func New() *Controller {
	return &Controller{
		validator: validator.NewCustomValidator(),
	}
}

func (ctrl *Controller) UpdatePanelConfig(c echo.Context) error {
	return nil
}
