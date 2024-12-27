package initialize

import (
	"api/internal/errors"
	"api/internal/models"
	"api/internal/repository"
	"api/pkg/config"
	"api/pkg/validator"
	"context"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"os"
)

type Controller struct {
	validator       validator.CustomValidatorInterface
	adminRepo       repository.AdminRepositoryInterface
	ocservGroupRepo repository.OcservGroupRepositoryInterface
}

func New() *Controller {
	return &Controller{
		validator:       validator.NewCustomValidator(),
		adminRepo:       repository.NewAdminRepository(),
		ocservGroupRepo: repository.NewOcservGroupRepository(),
	}
}

// CreateSuperUser Create Superuser account
//
// @Summary      Create Superuser
// @Description  Create Superuser in initializing step
// @Tags         init
// @Accept       json
// @Produce      json
// @Param        secret_key query string true "check secret key from file 'init_secret'"
// @Param        request    body  User   true "query params"
// @Success      200  {object}  nil
// @Router       /api/v1/init/admin [post]
func (ctrl *Controller) CreateSuperUser(c echo.Context) error {
	var user User
	if err := ctrl.validator.Validate(c, &user); err != nil {
		return errors.BadRequest(c, err.(error))
	}
	ctx := context.WithValue(c.Request().Context(), "username", user.Username)
	ctx = context.WithValue(ctx, "password", user.Password)
	err := ctrl.adminRepo.CreateSuperUser(ctx)
	go func() {
		err = os.Remove(config.GetApp().InitSecretFile)
		if err != nil {
			log.Println(err)
		}
	}()
	if err != nil {
		return errors.BadRequest(c, err)
	}
	return c.JSON(http.StatusCreated, nil)
}

func (ctrl *Controller) InitPanelConfig(c echo.Context) error {
	var data struct {
		GoogleCaptchaSecretKey string `json:"google_captcha_secret_key" validate:"omitempty"`
		GoogleCaptchaSiteKey   string `json:"google_captcha_site_key" validate:"omitempty"`
	}
	if err := ctrl.validator.Validate(c, &data); err != nil {
		return errors.BadRequest(c, err.(error))
	}
	ctx := context.WithValue(c.Request().Context(), "googleCaptchaSecretKey", data.GoogleCaptchaSecretKey)
	ctx = context.WithValue(ctx, "googleCaptchaSiteKey", data.GoogleCaptchaSiteKey)
	err := ctrl.adminRepo.CreateConfig(ctx)
	if err != nil {
		return errors.BadRequest(c, err)
	}
	return c.JSON(http.StatusCreated, nil)
}

func (ctrl *Controller) InitDefaultOcservGroup(c echo.Context) error {
	var data models.GroupConfig
	if err := ctrl.validator.Validate(c, &data); err != nil {
		return errors.BadRequest(c, err.(error))
	}
	ctx := context.WithValue(c.Request().Context(), "config", data)
	err := ctrl.ocservGroupRepo.UpdateDefaultGroup(ctx)
	if err != nil {
		return errors.BadRequest(c, err)
	}
	return c.JSON(http.StatusAccepted, nil)
}
