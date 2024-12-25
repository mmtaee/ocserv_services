package repository

import (
	"api/internal/handler"
	"api/internal/model"
	"api/pkg/database"
	"api/pkg/rabbitmq"
	"context"
	"errors"
	"gorm.io/gorm"
)

type OcservGroupRepository struct {
	db       *gorm.DB
	producer *rabbitmq.Producer
}

type OcservGroupRepositoryInterface interface {
	UpdateDefaultGroup(context.Context) error
	CreateGroup(context.Context) error
	UpdateGroup(context.Context) error
	DeleteGroup(context.Context) error
}

func NewOcservGroupRepository() *OcservGroupRepository {
	return &OcservGroupRepository{
		db:       database.Connection(),
		producer: rabbitmq.NewProducer(),
	}
}

func (o *OcservGroupRepository) UpdateDefaultGroup(ctx context.Context) error {
	config, ok := ctx.Value("config").(*model.GroupConfig)
	if !ok {
		config = &model.GroupConfig{}
	}
	ch := make(chan error, 1)
	go func() {
		processor := handler.NewProcessor(o.producer)
		event := handler.NewEvent(
			"group_default_update",
			map[string]interface{}{
				"name":   "group_default",
				"config": config,
			},
		)
		ch <- processor.ProcessEvent(event, "ocserv")
	}()
	return <-ch
}

func (o *OcservGroupRepository) CreateGroup(ctx context.Context) error {
	name := ctx.Value("name").(string)
	config, ok := ctx.Value("config").(map[string]string)
	if !ok {
		config = make(map[string]string)
	}
	ch := make(chan error, 1)
	go func() {
		processor := handler.NewProcessor(o.producer)
		event := handler.NewEvent(
			"group_create",
			map[string]interface{}{
				"name":   name,
				"config": config,
			},
		)
		ch <- processor.ProcessEvent(event, "ocserv")
	}()
	return <-ch
}

func (o *OcservGroupRepository) UpdateGroup(ctx context.Context) error {
	name := ctx.Value("name").(string)
	config, ok := ctx.Value("config").(map[string]string)
	if !ok {
		config = make(map[string]string)
	}
	action := "group_update"
	if name == "defaults" {
		action = "group_defaults_update"
	}
	ch := make(chan error, 1)
	go func() {
		processor := handler.NewProcessor(o.producer)
		event := handler.NewEvent(
			action,
			map[string]interface{}{
				"name":   name,
				"config": config,
			},
		)
		ch <- processor.ProcessEvent(event, "ocserv")
	}()
	return <-ch
}

func (o *OcservGroupRepository) DeleteGroup(ctx context.Context) error {
	ch := make(chan error, 1)
	name := ctx.Value("name").(string)
	if name == "defaults" {
		return errors.New("default group cannot be deleted")
	}
	go func() {
		processor := handler.NewProcessor(o.producer)
		event := handler.NewEvent("group_delete", map[string]string{
			"name": name,
		})
		ch <- processor.ProcessEvent(event, "ocserv")
	}()
	if err := <-ch; err != nil {
		return err
	}
	go func() {
		err := o.db.Table("users").Where("group = ?", name).Update("group", "defaults").Error
		ch <- err
	}()
	return <-ch
}
