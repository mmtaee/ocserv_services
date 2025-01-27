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
	Read(
		c context.Context,
		eventType string,
		conditions []string,
		args []interface{},
		order string,
		offset int,
		limit int,
		oldStateType interface{},
		newStateType interface{},
	) (*[]SchemaEvent, error)
}

func NewEventRepository(db *gorm.DB) *RepositoryEvent {
	return &RepositoryEvent{
		db: db,
	}
}

// Apply method for save Event on database
func (er *RepositoryEvent) Apply(c context.Context, schema *SchemaEvent) error {
	if err := schema.Validate(); err != nil {
		return err
	}
	event := schema.Serialize()
	return er.db.WithContext(c).Create(&event).Error
}

// Read method to fetch Events and Convert to []SchemaEvent
func (er *RepositoryEvent) Read(
	c context.Context,
	eventType string,
	conditions []string,
	args []interface{},
	order string,
	offset int,
	limit int,
	oldStateType interface{},
	newStateType interface{},
) (*[]SchemaEvent, error) {
	var events []Event
	var eventResults []SchemaEvent
	query := er.db.WithContext(c).
		Where("event_type = ?", eventType).
		Offset(offset).
		Limit(limit).
		Order(order)
	if conditions != nil {
		query = query.Where(conditions, args...)
	}
	err := query.Find(&events).Error
	if err != nil {
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
