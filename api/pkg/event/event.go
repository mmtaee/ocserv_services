package event

import (
	"encoding/json"
	"errors"
	"github.com/mmtaee/go-oc-utils/logger"
	"reflect"
	"time"
)

// Event struct database model
type Event struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	EventType string    `json:"event_type" gorm:"type:varchar(32);not null"`
	ModelName string    `json:"model_name" gorm:"type:varchar(32);not null"`
	ModelUID  string    `json:"model_uid" gorm:"type:varchar(32)"`
	UserUID   string    `json:"user_uid" gorm:"type:varchar(32)"`
	OldState  string    `json:"old_state" gorm:"type:text"`
	NewState  string    `json:"new_state" gorm:"type:text"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
}

// SchemaEvent struct request schema
type SchemaEvent struct {
	ID        uint        `json:"id"`
	ModelName string      `json:"model_name" `
	ModelUID  string      `json:"model_uid"`
	EventType string      `json:"event_type"`
	UserUID   string      `json:"user_uid"`
	OldState  interface{} `json:"old_state"`
	NewState  interface{} `json:"new_state"`
}

// toJSON method convert old and new state to json
func toJSON(data interface{}) (string, error) {
	if data == nil {
		return "", nil
	}
	if str, ok := data.(string); ok {
		return str, nil
	}
	res, err := json.Marshal(data)
	if err != nil {
		return "{}", err
	}
	return string(res), nil
}

// Serialize method for create Event type from SchemaEvent
func (e *SchemaEvent) Serialize() *Event {
	var err error

	event := &Event{
		ModelName: e.ModelName,
		ModelUID:  e.ModelUID,
		EventType: e.EventType,
		UserUID:   e.UserUID,
	}

	event.OldState, err = toJSON(e.OldState)
	if err != nil {
		logger.Logf(logger.ERROR, "Failed to serialize old state: %v", err)
	}

	event.NewState, err = toJSON(e.NewState)
	if err != nil {
		logger.Logf(logger.ERROR, "Failed to serialize new state: %v", err)
	}
	return event
}

// Deserialize method for convert Event to SchemaEvent
func (e *Event) Deserialize(oldStateType, newStateType interface{}) (*SchemaEvent, error) {
	if e.OldState != "" {
		if reflect.TypeOf(oldStateType).Kind() == reflect.String {
			oldStateType = e.OldState
		} else {
			err := json.Unmarshal([]byte(e.OldState), oldStateType)
			logger.Logf(logger.ERROR, "Failed to deserialize old state: %v", err)
		}
	} else {
		oldStateType = nil
	}
	if e.NewState != "" {
		if reflect.TypeOf(newStateType).Kind() == reflect.String {
			newStateType = e.NewState
		} else {
			err := json.Unmarshal([]byte(e.NewState), newStateType)
			if err != nil {
				logger.Logf(logger.ERROR, "Failed to deserialize new state: %v", err)
			}
		}
	} else {
		newStateType = nil
	}

	return &SchemaEvent{
		ID:        e.ID,
		ModelName: e.ModelName,
		ModelUID:  e.ModelUID,
		EventType: e.EventType,
		UserUID:   e.UserUID,
		OldState:  oldStateType,
		NewState:  newStateType,
	}, nil
}

// Validate SchemaEvent
func (e *SchemaEvent) Validate() error {
	if e.ModelName == "" || e.UserUID == "" {
		return errors.New("missing required fields in Event")
	}
	return nil
}
