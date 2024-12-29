package repository

import (
	"api/internal/models"
	"api/pkg/database"
	"api/pkg/password"
	"context"
	"gorm.io/gorm"
)

type AdminRepository struct {
	db *gorm.DB
}

type AdminRepositoryInterface interface {
	CreateSuperUser(context.Context, string, string) error
}

func NewAdminRepository() *AdminRepository {
	return &AdminRepository{
		db: database.Connection(),
	}
}

func (a *AdminRepository) CreateSuperUser(c context.Context, username, passwd string) error {
	passwordHash := password.Create(passwd)
	return a.db.WithContext(c).Create(&models.User{
		Username: username,
		Password: passwordHash,
	}).Error

}
