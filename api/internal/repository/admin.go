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
	CreateSuperUser(context.Context) error
}

func NewAdminRepository() *AdminRepository {
	return &AdminRepository{
		db: database.Connection(),
	}
}

func (a *AdminRepository) CreateSuperUser(c context.Context) error {
	ch := make(chan error, 1)
	passwordString := c.Value("password").(string)
	go func() {
		passwordHash := password.Create(passwordString)
		ch <- a.db.WithContext(c).Create(&models.User{
			Username: c.Value("username").(string),
			Password: passwordHash,
		}).Error
	}()
	return <-ch
}
