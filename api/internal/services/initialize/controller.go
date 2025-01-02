package initialize

import (
	"api/internal/models"
	"api/internal/repository"
	"api/pkg/config"
	"api/pkg/utils"
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

// CheckSecretKey check and validate secret key
//
// @Summary      Check and validate secret key
// @Description  Check and validate secret key for continue initialing app
// @Tags         Initial
// @Accept       json
// @Produce      json
// @Param        secret query string true "check secret key from file 'init_secret'"
// @Success      200  {object} nil
// @Failure      400 {object} utils.ErrorResponse
// @Router       /api/v1/init/check [get]
func (ctrl *Controller) CheckSecretKey(c echo.Context) error {
	if err := checkSecret(c.QueryParam("secret")); err != nil {
		return utils.BadRequest(c, err)
	}
	return c.JSON(http.StatusOK, nil)
}

// CreateSuperUser Create Superuser account
//
// @Summary      Create Superuser
// @Description  Create Superuser in initializing step
// @Tags         Initial
// @Accept       json
// @Produce      json
// @Param        secret query string true "check secret key from file 'init_secret'"
// @Param        request body  CreateAdminUserRequest true "admin user body data"
// @Success      200  {object} nil
// @Failure      400 {object} utils.ErrorResponse
// @Router       /api/v1/init/admin [post]
func (ctrl *Controller) CreateSuperUser(c echo.Context) error {
	if err := checkSecret(c.QueryParam("secret")); err != nil {
		return utils.BadRequest(c, err)
	}
	var user CreateAdminUserRequest
	if err := ctrl.validator.Validate(c, &user); err != nil {
		return utils.BadRequest(c, err.(error))
	}
	err := ctrl.userRepo.Admin.CreateSuperUser(c.Request().Context(), user.Username, user.Password)
	if err != nil {
		return utils.BadRequest(c, err)
	}
	return c.JSON(http.StatusCreated, nil)
}

// PanelConfig Create Panel Config
//
// @Summary      Create Panel Config
// @Description  Create Panel Config initializing step
// @Tags         Initial
// @Accept       json
// @Produce      json
// @Param        init_secret query string true "check secret key from file 'init_secret'"
// @Param        request    body  CreateSiteConfigRequest   true "site config data"
// @Success      201  {object}  nil
// @Failure      400 {object} utils.ErrorResponse
// @Router       /api/v1/init/config [post]
func (ctrl *Controller) PanelConfig(c echo.Context) error {
	if err := checkSecret(c.QueryParam("secret")); err != nil {
		return utils.BadRequest(c, err)
	}
	var data CreateSiteConfigRequest
	if err := ctrl.validator.Validate(c, &data); err != nil {
		return utils.BadRequest(c, err.(error))
	}
	panelConfig := models.PanelConfig{
		GoogleCaptchaSecretKey: data.GoogleCaptchaSecretKey,
		GoogleCaptchaSiteKey:   data.GoogleCaptchaSiteKey,
	}
	err := ctrl.panelRepo.CreateConfig(c.Request().Context(), panelConfig)
	if err != nil {
		return utils.BadRequest(c, err)
	}
	go func() {
		err = os.Remove(config.GetApp().InitSecretFile)
		if err != nil {
			log.Println(err)
		}
	}()
	return c.JSON(http.StatusCreated, nil)
}
