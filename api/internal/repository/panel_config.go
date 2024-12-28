package repository

import (
	"api/internal/models"
	"api/pkg/database"
	"context"
	"gorm.io/gorm"
)

type PanelConfigRepository struct {
	db *gorm.DB
}

type PanelConfigRepositoryInterface interface {
	CreateConfig(context.Context) error
}

func NewPanelConfigRepository() PanelConfigRepositoryInterface {
	return &PanelConfigRepository{
		db: database.Connection(),
	}
}

func (p *PanelConfigRepository) CreateConfig(c context.Context) error {
	ch := make(chan error, 1)
	go func() {
		ch <- p.db.WithContext(c).Create(&models.PanelConfig{
			GoogleCaptchaSecretKey: c.Value("googleCaptchaSecretKey").(string),
			GoogleCaptchaSiteKey:   c.Value("googleCaptchaSiteKey").(string),
		}).Error
	}()
	return <-ch
}
