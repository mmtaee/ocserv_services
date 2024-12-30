package staff

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

func (ctrl *Controller) Staffs(c echo.Context) error {
	return nil
}

func (ctrl *Controller) Staff(c echo.Context) error {
	// TODO: detail with permission
	return nil
}

func (ctrl *Controller) CreateStaff(c echo.Context) error {
	// TODO: create Staff + Permission
	return nil
}

func (ctrl *Controller) UpdateStaff(c echo.Context) error {
	// TODO: update Staff + Permission
	return nil
}

func (ctrl *Controller) DeleteStaff(c echo.Context) error {
	return nil
}
