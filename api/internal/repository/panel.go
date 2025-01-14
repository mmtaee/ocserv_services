package repository

import (
	"api/internal/models"
	"context"
	"github.com/mmtaee/go-oc-utils/database"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PanelConfigRepository struct {
	db *gorm.DB
}

type PanelConfigRepositoryInterface interface {
	CreateConfig(c context.Context, config models.PanelConfig) error
	UpdateConfig(c context.Context, siteKey, secretKet string) error
	GetConfig(c context.Context) (*models.PanelConfig, error)
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

func (p *PanelConfigRepository) GetConfig(c context.Context) (*models.PanelConfig, error) {
	config := &models.PanelConfig{}
	err := p.db.WithContext(c).Model(models.PanelConfig{}).First(&config).Error
	if err != nil {
		return nil, err
	}
	return config, nil
}
