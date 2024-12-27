package repository

import (
	"api/pkg/database"
	"gorm.io/gorm"
)

type OcservUserRepository struct {
	db *gorm.DB
}

type OcservUserRepositoryInterface interface{}

func NewOcservUserRepository() *OcservUserRepository {
	return &OcservUserRepository{
		db: database.Connection(),
	}
}
