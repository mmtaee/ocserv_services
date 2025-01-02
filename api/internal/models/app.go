package models

import (
	"api/pkg/database"
	"errors"
	"gorm.io/gorm"
)

type PanelConfig struct {
	ID                     uint   `json:"-" gorm:"primaryKey;autoIncrement"`
	GoogleCaptchaSecretKey string `json:"google_captcha_secret" gorm:"type:text"`
	GoogleCaptchaSiteKey   string `json:"google_captcha_site_key" gorm:"type:text"`
}

func (p *PanelConfig) BeforeCreate(tx *gorm.DB) error {
	ch := make(chan error, 1)
	go func() {
		var config PanelConfig
		db := database.Connection()
		err := db.Table("panel_configs").First(&config).Error
		if err != nil && config.ID == 0 {
			ch <- nil
		}
		ch <- errors.New("panel configs already exist")
	}()
	return <-ch
}
