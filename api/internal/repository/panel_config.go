package repository

import (
	"api/internal/models"
	"api/pkg/database"
	"context"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PanelConfigRepository struct {
	db *gorm.DB
}

type PanelConfigRepositoryInterface interface {
	CreateConfig(context.Context, models.PanelConfig) error
	UpdateConfig(context.Context, string, string) error
}

func NewPanelConfigRepository() PanelConfigRepositoryInterface {
	return &PanelConfigRepository{
		db: database.Connection(),
	}
}

func (p *PanelConfigRepository) CreateConfig(c context.Context, config models.PanelConfig) error {
	return p.db.WithContext(c).Create(&config).Error
}

func (p *PanelConfigRepository) UpdateConfig(c context.Context, siteKey, secretKet string) error {
	config := &models.PanelConfig{}
	return p.db.WithContext(c).Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&config).Error; err != nil {
			return err
		}
		config.GoogleCaptchaSiteKey = siteKey
		config.GoogleCaptchaSecretKey = secretKet
		return tx.WithContext(c).Save(&config).Error
	})
}
