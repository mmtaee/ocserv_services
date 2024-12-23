package repository

import (
	"api/pkg/postgres"
	"api/pkg/rabbitmq"
	"xorm.io/xorm"
)

type OcservUserRepository struct {
	db       *xorm.Engine
	producer *rabbitmq.Producer
}

type OcservUserRepositoryInterface interface{}

func NewOcservUserRepository() *OcservGroupRepository {
	return &OcservGroupRepository{
		db:       postgres.GetEngine(),
		producer: rabbitmq.NewProducer(),
	}
}
