package repository

import (
	"api/internal/models"
	"api/pkg/database"
	"context"
	"gorm.io/gorm"
)

type AdminRepository struct {
	db *gorm.DB
}

type AdminRepositoryInterface interface {
	CreateSuperUser(context.Context) error
	CreateConfig(context.Context) error
}

func NewAdminRepository() *AdminRepository {
	return &AdminRepository{
		db: database.Connection(),
	}
}

func (r *AdminRepository) CreateSuperUser(c context.Context) error {
	ch := make(chan error, 1)
	go func() {
		ch <- r.db.WithContext(c).Create(&models.Admin{
			Username: c.Value("username").(string),
			Password: c.Value("password").(string),
		}).Error
	}()
	return <-ch
}

func (r *AdminRepository) CreateConfig(c context.Context) error {
	ch := make(chan error, 1)
	go func() {
		ch <- r.db.WithContext(c).Create(&models.PanelConfig{
			GoogleCaptchaSecretKey: c.Value("googleCaptchaSecretKey").(string),
			GoogleCaptchaSiteKey:   c.Value("googleCaptchaSiteKey").(string),
		}).Error
	}()
	return <-ch
}
