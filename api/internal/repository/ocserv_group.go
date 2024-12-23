package repository

import (
	"api/internal/handler"
	"api/internal/model"
	"api/pkg/postgres"
	"api/pkg/rabbitmq"
	"context"
	"errors"
	"xorm.io/xorm"
)

type OcservGroupRepository struct {
	db       *xorm.Engine
	producer *rabbitmq.Producer
}

type OcservGroupRepositoryInterface interface {
	CreateGroup(context.Context) error
	UpdateGroup(context.Context) error
	DeleteGroup(context.Context) error
}

func NewOcservGroupRepository() *OcservGroupRepository {
	return &OcservGroupRepository{
		db:       postgres.GetEngine(),
		producer: rabbitmq.NewProducer(),
	}
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
		_, err := o.db.Table(&model.User{}).Where("group = ?", name).Update(map[string]interface{}{
			"group": "defaults",
		})
		ch <- err
	}()
	return <-ch
}
