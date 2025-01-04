package panel

import (
	"api/internal/models"
	"api/internal/repository"
	"api/pkg/utils"
	"api/pkg/validator"
	"github.com/labstack/echo/v4"
	"net/http"
)

type Controller struct {
	validator validator.CustomValidatorInterface
	panelRepo repository.PanelConfigRepositoryInterface
}

func New() *Controller {
	return &Controller{
		validator: validator.NewCustomValidator(),
		panelRepo: repository.NewPanelConfigRepository(),
	}
}

// CreatePanelConfig Create Panel Config
//
// @Summary      Create Panel Config
// @Description  Create Panel Config initializing step
// @Tags         Panel
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer TOKEN"
// @Param        init_secret query string true "check secret key from file 'init_secret'"
// @Param        request    body  CreateSiteConfigRequest   true "site config data"
// @Success      201  {object}  nil
// @Failure      400 {object} utils.ErrorResponse
// @Failure      401 {object} middlewares.Unauthorized
// @Router       /api/v1/panel/config [post]
func (ctrl *Controller) CreatePanelConfig(c echo.Context) error {
	var data CreateSiteConfigRequest
	if err := ctrl.validator.Validate(c, &data); err != nil {
		return utils.BadRequest(c, err.(error))
	}
	panelConfig := models.PanelConfig{
		Init:                   true,
		GoogleCaptchaSecretKey: data.GoogleCaptchaSecretKey,
		GoogleCaptchaSiteKey:   data.GoogleCaptchaSiteKey,
	}
	err := ctrl.panelRepo.CreateConfig(c.Request().Context(), panelConfig)
	if err != nil {
		return utils.BadRequest(c, err)
	}
	return c.JSON(http.StatusCreated, nil)
}

// UpdatePanelConfig Update Panel Config
//
// @Summary      Update Panel Config
// @Description  Update Panel Config in app after login step
// @Tags         Panel
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer TOKEN"
// @Param        request    body  UpdateSiteConfigRequest   true "site config data"
// @Success      200  {object}  nil
// @Failure      400 {object} utils.ErrorResponse
// @Failure      401 {object} middlewares.Unauthorized
// @Router       /api/v1/panel/config [put]
func (ctrl *Controller) UpdatePanelConfig(c echo.Context) error {
	var data UpdateSiteConfigRequest
	if err := ctrl.validator.Validate(c, &data); err != nil {
		return utils.BadRequest(c, err.(error))
	}
	err := ctrl.panelRepo.UpdateConfig(c.Request().Context(), data.GoogleCaptchaSiteKey, data.GoogleCaptchaSecretKey)
	if err != nil {
		return utils.BadRequest(c, err)
	}
	return c.JSON(http.StatusOK, nil)
}

// GetPanelInitConfig Get Panel Config
//
// @Summary      Get Panel Config with Init
// @Description  Get Panel Config to discover init data
// @Tags         Panel
// @Accept       json
// @Produce      json
// @Success      200  {object}  GetPanelConfigResponse
// @Failure      400 {object} utils.ErrorResponse
// @Router       /api/v1/panel/config/init [get]
func (ctrl *Controller) GetPanelInitConfig(c echo.Context) error {
	config, err := ctrl.panelRepo.GetConfig(c.Request().Context())
	if err != nil {
		return utils.BadRequest(c, err)
	}
	return c.JSON(http.StatusOK, GetPanelConfigResponse{
		Init:                 config.Init,
		GoogleCaptchaSiteKey: config.GoogleCaptchaSiteKey,
	})
}

// GetPanelConfig Get Panel Config
//
// @Summary      Get Panel Config
// @Description  Get Panel Config
// @Tags         Panel
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer TOKEN"
// @Success      200  {object}  GetFullPanelConfigResponse
// @Failure      400 {object} utils.ErrorResponse
// @Failure      401 {object} middlewares.Unauthorized
// @Router       /api/v1/panel/config [get]
func (ctrl *Controller) GetPanelConfig(c echo.Context) error {
	config, err := ctrl.panelRepo.GetConfig(c.Request().Context())
	if err != nil {
		return utils.BadRequest(c, err)
	}
	return c.JSON(http.StatusOK, GetFullPanelConfigResponse{
		GoogleCaptchaSiteKey:   config.GoogleCaptchaSiteKey,
		GoogleCaptchaSecretKey: config.GoogleCaptchaSecretKey,
	})
}
