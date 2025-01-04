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
	CreateSuperUser(c context.Context, username, passwd string) (*models.User, error)
}

func NewAdminRepository() *AdminRepository {
	return &AdminRepository{
		db: database.Connection(),
	}
}

func (a *AdminRepository) CreateSuperUser(c context.Context, username, passwd string) (*models.User, error) {
	user := models.User{
		Username: username,
		Password: password.Create(passwd),
	}
	err := a.db.WithContext(c).Create(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
