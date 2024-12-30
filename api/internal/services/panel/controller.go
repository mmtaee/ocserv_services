package panel

import (
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

// UpdatePanelConfig Update Panel Config
//
// @Summary      Update Panel Config
// @Description  Update Panel Config in app after login step
// @Tags         Panel
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer TOKEN"
// @Failure      401 {object} middlewares.Unauthorized
// @Param        request    body  UpdateSiteConfigRequest   true "site config data"
// @Success      200  {object}  nil
// @Failure      400 {object} utils.ErrorResponse
// @Router       /services/v1/panel/config [put]
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
