package repository

import (
	"api/pkg/database"
	"api/pkg/rabbitmq"
	"gorm.io/gorm"
)

type OcservUserRepository struct {
	db       *gorm.DB
	producer *rabbitmq.Producer
}

type OcservUserRepositoryInterface interface{}

func NewOcservUserRepository() *OcservGroupRepository {
	return &OcservGroupRepository{
		db:       database.Connection(),
		producer: rabbitmq.NewProducer(),
	}
}
