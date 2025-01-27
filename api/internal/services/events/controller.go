package events

import (
	"api/internal/repository"
	"api/pkg/utils"
	"errors"
	"github.com/labstack/echo/v4"
	"net/http"
	"slices"
)

type Controller struct {
	validator utils.CustomValidatorInterface
	eventRepo repository.EventsRepositoryInterface
}

func New() *Controller {
	return &Controller{
		validator: utils.NewCustomValidator(),
		eventRepo: repository.NewEventRepository(),
	}
}

var EventModels = []string{
	"create_staff",
	"create_staff_permission",
	"update_staff_permission",
	"update_staff_password",
	"delete_staff",
	"update_panel_config",
	"update_oc_default_group",
	"create_oc_group",
	"update_oc_group",
	"delete_oc_group",
	"create_oc_user",
	"update_oc_user",
	"lock_oc_user",
	"unlock_oc_user",
	"disconnect_oc_user",
	"delete_oc_user",
}

// Events List of events
//
// @Summary      List of events
// @Description  List of events
// @Tags         Events
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer TOKEN"
// @Param 		 page query int false "Page number, starting from 1" minimum(1)
// @Param 		 pager query string false "Field to order by"
// @Param 		 sort query string false "Sort order, either ASC or DESC" Enums(ASC, DESC)
// @Param 		 model_name query string false "event model name"
// @Param 		 event_type query string false "event type"
// @Param 		 user_id query string false "id of user that does this event"
// @Param 		 date_start query string false "event date create from"
// @Param 		 date_end query string false "event date create to"
// @Success      200 {object} []events.SchemaEvent
// @Failure      400 {object} utils.ErrorResponse
// @Failure      401 {object} middlewares.Unauthorized
// @Router       /api/v1/events/:event_type [get]
func (ctrl *Controller) Events(c echo.Context) error {
	eventType := c.Param("event_type")
	if !slices.Contains(EventModels, eventType) {
		return utils.BadRequest(c, errors.New("invalid event name"))
	}
	pageData := utils.NewPaginationRequest()
	var filters repository.EventFilterRequest
	if err := ctrl.validator.Validate(c, &pageData); err != nil {
		return utils.BadRequest(c, err)
	}
	if err := ctrl.validator.Validate(c, &filters); err != nil {
		return utils.BadRequest(c, err)
	}
	events, err := ctrl.eventRepo.Events(c.Request().Context(), eventType, &pageData, &filters)
	if err != nil {
		return utils.BadRequest(c, err)
	}
	return c.JSON(http.StatusOK, events)
}
