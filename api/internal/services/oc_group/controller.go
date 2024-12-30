package oc_group

import (
	"api/internal/cu_errors"
	"api/internal/models"
	"api/internal/repository"
	"api/pkg/validator"
	"github.com/labstack/echo/v4"
	"net/http"
)

type Controller struct {
	validator       validator.CustomValidatorInterface
	ocservGroupRepo repository.OcservGroupRepositoryInterface
}

func New() *Controller {
	return &Controller{
		validator:       validator.NewCustomValidator(),
		ocservGroupRepo: repository.NewOcservGroupRepository(),
	}
}

// UpdateDefaultOcservGroup Create Superuser account
//
// @Summary      Update Ocserv Defaults Group
// @Description  Update Ocserv Defaults Group initializing step
// @Tags         Ocserv Group
// @Accept       json
// @Produce      json
// @Failure      401 {object} middlewares.Unauthorized
// @Param        request body  models.OcGroupConfig true "oc group default config"
// @Success      200  {object}  nil
// @Failure      400 {object} cu_errors.ErrorResponse
// @Router       /api/v1/ocserv/group [post]
func (ctrl *Controller) UpdateDefaultOcservGroup(c echo.Context) error {
	var data models.OcGroupConfig
	if err := ctrl.validator.Validate(c, &data); err != nil {
		return cu_errors.BadRequest(c, err.(error))
	}
	err := ctrl.ocservGroupRepo.UpdateDefaultGroup(c.Request().Context(), data)
	if err != nil {
		return cu_errors.BadRequest(c, err)
	}
	return c.JSON(http.StatusAccepted, nil)
}
