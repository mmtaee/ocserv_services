package repository

import (
	"api/pkg/postgres"
	"api/pkg/rabbitmq"
	"xorm.io/xorm"
)

type AdminRepository struct {
	db       *xorm.Engine
	producer *rabbitmq.Producer
}

type AdminRepositoryInterface interface {
}

func NewAdminRepository() *AdminRepository {
	return &AdminRepository{
		db:       postgres.GetEngine(),
		producer: rabbitmq.NewProducer(),
	}
}
