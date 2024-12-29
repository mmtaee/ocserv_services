package initialize

import (
	"api/internal/errors"
	"api/internal/models"
	"api/internal/repository"
	"api/pkg/config"
	"api/pkg/validator"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"os"
)

type Controller struct {
	validator       validator.CustomValidatorInterface
	userRepo        *repository.UserRepository
	panelRepo       repository.PanelConfigRepositoryInterface
	ocservGroupRepo repository.OcservGroupRepositoryInterface
}

func New() *Controller {
	return &Controller{
		validator:       validator.NewCustomValidator(),
		userRepo:        repository.NewUserRepository(),
		panelRepo:       repository.NewPanelConfigRepository(),
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
// @Param        request    body  CreateAdminUserRequest   true "query params"
// @Success      200  {object} nil
// @Failure      400 {object} errors.ErrorResponse
// @Router       /api/v1/init/admin [post]
func (ctrl *Controller) CreateSuperUser(c echo.Context) error {
	var user CreateAdminUserRequest
	if err := ctrl.validator.Validate(c, &user); err != nil {
		return errors.BadRequest(c, err.(error))
	}
	err := ctrl.userRepo.Admin.CreateSuperUser(c.Request().Context(), user.Username, user.Password)
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

// PanelConfig Create Superuser account
//
// @Summary      Create Panel Config
// @Description  Create Panel Config initializing step
// @Tags         init
// @Accept       json
// @Produce      json
// @Param        secret_key query string true "check secret key from file 'init_secret'"
// @Param        request    body  CreateSiteConfigRequest   true "query params"
// @Success      200  {object}  nil
// @Failure      400 {object} errors.ErrorResponse
// @Router       /api/v1/init/config [post]
func (ctrl *Controller) PanelConfig(c echo.Context) error {
	var data CreateSiteConfigRequest
	if err := ctrl.validator.Validate(c, &data); err != nil {
		return errors.BadRequest(c, err.(error))
	}
	panelConfig := models.PanelConfig{
		GoogleCaptchaSecretKey: data.GoogleCaptchaSecretKey,
		GoogleCaptchaSiteKey:   data.GoogleCaptchaSiteKey,
	}
	err := ctrl.panelRepo.CreateConfig(c.Request().Context(), panelConfig)
	if err != nil {
		return errors.BadRequest(c, err)
	}
	return c.JSON(http.StatusCreated, nil)
}

// DefaultOcservGroup Create Superuser account
//
// @Summary      Update Ocserv Group
// @Description  Update Ocserv Defaults Group initializing step
// @Tags         init
// @Accept       json
// @Produce      json
// @Param        secret_key query string true "check secret key from file 'init_secret'"
// @Param        request    body  models.OcGroupConfig   true "query params"
// @Success      200  {object}  nil
// @Failure      400 {object} errors.ErrorResponse
// @Router       /api/v1/init/group [post]
func (ctrl *Controller) DefaultOcservGroup(c echo.Context) error {
	var data models.OcGroupConfig
	if err := ctrl.validator.Validate(c, &data); err != nil {
		return errors.BadRequest(c, err.(error))
	}
	err := ctrl.ocservGroupRepo.UpdateDefaultGroup(c.Request().Context(), data)
	if err != nil {
		return errors.BadRequest(c, err)
	}
	return c.JSON(http.StatusAccepted, nil)
}
