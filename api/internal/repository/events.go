package repository

import (
	"api/pkg/event"
	"api/pkg/utils"
	"context"
	"errors"
	"github.com/mmtaee/go-oc-utils/database"
	"github.com/mmtaee/go-oc-utils/handler/ocgroup"
	"github.com/mmtaee/go-oc-utils/models"
)

type EventRepository struct {
	eventRepo event.RepositoryEventInterface
}

type EventsRepositoryInterface interface {
	Events(c context.Context, eventType string, page *utils.RequestPagination, filters *EventFilterRequest) (*[]event.SchemaEvent, error)
}

type EventFilterRequest struct {
	ModelName *string `query:"model_name"`
	UserId    *string `query:"user_id"`
	DateStart *string `query:"date_start"`
	DateEnd   *string `query:"date_end"`
}

func NewEventRepository() *EventRepository {
	return &EventRepository{eventRepo: event.NewEventRepository(database.Connection())}
}

func (e *EventRepository) Events(
	c context.Context,
	eventType string,
	page *utils.RequestPagination,
	filters *EventFilterRequest,
) (*[]event.SchemaEvent, error) {
	var oldStateType interface{}
	var newStateType interface{}
	switch eventType {
	case "create_staff":
		oldStateType = &models.User{}
		newStateType = &models.User{}
	case "create_staff_permission":
		oldStateType = &models.UserPermission{}
		newStateType = &models.UserPermission{}
	case "update_staff_permission":
		oldStateType = &models.UserPermission{}
		newStateType = &models.UserPermission{}
	case "update_staff_password":
		oldStateType = nil
		newStateType = nil
	case "delete_staff":
		oldStateType = nil
		newStateType = nil
	case "update_panel_config":
		oldStateType = &models.PanelConfig{}
		newStateType = &models.PanelConfig{}
	case "update_oc_default_group":
		oldStateType = &ocgroup.OcservGroupConfig{}
		newStateType = &ocgroup.OcservGroupConfig{}
	case "create_oc_group":
		oldStateType = &ocgroup.OcservGroup{}
		newStateType = &ocgroup.OcservGroup{}
	case "update_oc_group":
		oldStateType = &ocgroup.OcservGroupConfig{}
		newStateType = &ocgroup.OcservGroupConfig{}
	case "delete_oc_group":
		oldStateType = nil
		newStateType = nil

	case "create_oc_user":
		oldStateType = &models.OcUser{}
		newStateType = &models.OcUser{}
	case "update_oc_user":
		oldStateType = &models.OcUser{}
		newStateType = &models.OcUser{}
	case "lock_oc_user":
		oldStateType = ""
		newStateType = ""
	case "unlock_oc_user":
		oldStateType = ""
		newStateType = ""
	case "disconnect_oc_user":
		oldStateType = nil
		newStateType = nil
	case "delete_oc_user":
		oldStateType = nil
		newStateType = nil
	default:
		return nil, errors.New("not found")
	}
	var conditions []string
	var args []interface{}
	if filters.ModelName != nil {
		conditions = append(conditions, "model_name = ?")
		args = append(args, *filters.ModelName)
	}
	if filters.UserId != nil {
		conditions = append(conditions, "user_id = ?")
		args = append(args, *filters.UserId)
	}
	if filters.DateStart != nil {
		conditions = append(conditions, "created_at >= ?")
		args = append(args, *filters.DateStart)
	}
	if filters.DateEnd != nil {
		conditions = append(conditions, "created_at <= ?")
		args = append(args, *filters.DateEnd)
	}

	offset := (page.Page - 1) * page.PageSize
	eventList, err := e.eventRepo.Read(c, eventType, conditions, args, page.Order, offset, page.PageSize, oldStateType, newStateType)
	if err != nil {
		return nil, err
	}
	return eventList, nil
}
