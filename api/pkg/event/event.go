package event

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mmtaee/go-oc-utils/logger"
	"time"
)

type Event struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	ModelName string    `json:"model_name" gorm:"type:varchar(32);not null"`
	ModelUID  string    `json:"model_uid" gorm:"type:varchar(32);not null"`
	EventType string    `json:"event_type" gorm:"type:varchar(32);not null"`
	UserUID   string    `json:"user_uid" gorm:"type:varchar(32);not null"`
	OldState  string    `json:"old_state" gorm:"type:text"`
	NewState  string    `json:"new_state" gorm:"type:text"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
}

type SchemaEvent struct {
	ID        uint        `json:"id"`
	ModelName string      `json:"model_name" `
	ModelUID  string      `json:"model_uid"`
	EventType string      `json:"event_type"`
	UserUID   string      `json:"user_uid"`
	OldState  interface{} `json:"old_state"`
	NewState  interface{} `json:"new_state"`
}

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

func (e *Event) Deserialize(oldStateType, newStateType interface{}) (*SchemaEvent, error) {
	err := json.Unmarshal([]byte(e.OldState), oldStateType)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize OldState: %w", err)
	}

	err = json.Unmarshal([]byte(e.NewState), newStateType)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize NewState: %w", err)
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

func (e *SchemaEvent) Validate() error {
	if e.ModelName == "" || e.UserUID == "" {
		return errors.New("missing required fields in Event")
	}
	return nil
}
