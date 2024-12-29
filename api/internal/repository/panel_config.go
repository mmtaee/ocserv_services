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
	CreateConfig(context.Context, models.PanelConfig) error
}

func NewPanelConfigRepository() PanelConfigRepositoryInterface {
	return &PanelConfigRepository{
		db: database.Connection(),
	}
}

func (p *PanelConfigRepository) CreateConfig(c context.Context, config models.PanelConfig) error {
	return p.db.WithContext(c).Create(&config).Error
}
