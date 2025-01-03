package ocUser

import (
	"api/internal/repository"
	"api/pkg/database"
	"api/pkg/validator"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Controller struct {
	db             *gorm.DB
	validator      validator.CustomValidatorInterface
	ocservUserRepo repository.OcservUserRepositoryInterface
}

func New() *Controller {
	return &Controller{
		db:             database.Connection(),
		validator:      validator.NewCustomValidator(),
		ocservUserRepo: repository.NewOcservUserRepository(),
	}
}

func (ctrl *Controller) Users(c echo.Context) error {
	return nil
}

func (ctrl *Controller) User(c echo.Context) error {
	return nil
}

func (ctrl *Controller) Create(c echo.Context) error {
	return nil
}

func (ctrl *Controller) Update(c echo.Context) error {
	return nil
}

func (ctrl *Controller) Lock(c echo.Context) error {
	// TODO: get lock true or false from body
	return nil
}

func (ctrl *Controller) Disconnect(c echo.Context) error {
	return nil
}

func (ctrl *Controller) Delete(c echo.Context) error {
	return nil
}

func (ctrl *Controller) Statistics(c echo.Context) error {
	return nil
}

func (ctrl *Controller) Activity(c echo.Context) error {
	return nil
}
