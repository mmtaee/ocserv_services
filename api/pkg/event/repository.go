package event

import (
	"context"
	"gorm.io/gorm"
)

type RepositoryEvent struct {
	db *gorm.DB
}

type RepositoryEventInterface interface {
	Apply(c context.Context, schema *SchemaEvent) error
	Read(c context.Context, eventModel string, oldStateType interface{}, newStateType interface{}) (*[]SchemaEvent, error)
}

func NewEventRepository(db *gorm.DB) *RepositoryEvent {
	return &RepositoryEvent{
		db: db,
	}
}

func (er *RepositoryEvent) Apply(c context.Context, schema *SchemaEvent) error {
	if err := schema.Validate(); err != nil {
		return err
	}
	event := schema.Serialize()
	return er.db.WithContext(c).Create(&event).Error
}

func (er *RepositoryEvent) Read(c context.Context, modelName string, oldStateType, newStateType interface{}) (*[]SchemaEvent, error) {
	var events []Event
	var eventResults []SchemaEvent
	if err := er.db.WithContext(c).Where("model_name = ?", modelName).Find(&events).Error; err != nil {
		return nil, err
	}

	for _, event := range events {
		result, err := event.Deserialize(oldStateType, newStateType)
		if err != nil {
			return nil, err
		}
		eventResults = append(eventResults, *result)
	}
	return &eventResults, nil
}
