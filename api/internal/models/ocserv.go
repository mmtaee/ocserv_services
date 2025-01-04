package models

import (
	"encoding/json"
	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
	"time"
)

const (
	Free int32 = iota
	Monthly
	Totally
)

type OcUser struct {
	ID          uint      `json:"-" gorm:"primaryKey;autoIncrement"`
	UID         string    `json:"uid" gorm:"type:varchar(26);not null;unique"`
	Group       string    `json:"group" gorm:"type:varchar(16);default('defaults')"`
	Username    string    `json:"username" gorm:"type:varchar(16);not null;unique"`
	Password    string    `json:"password" gorm:"type:varchar(16);not null"`
	IsLocked    bool      `json:"isLocked" gorm:"default(false)"`
	CreatedAt   time.Time `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updatedAt" gorm:"autoUpdateTime"`
	ExpireAt    time.Time `json:"expire_at"`
	TrafficType int32     `json:"trafficType" gorm:"not null;default(0)"`
	TrafficSize int32     `json:"trafficSize" gorm:"not null;default(10)"`
	Rx          float64   `json:"rx" gorm:"not null;default(0.00)"`
	Tx          float64   `json:"tx" gorm:"not null;default(0.00)"`
	Description string    `json:"description" gorm:"type:text"`
	IsOnline    bool      `json:"is_online" gorm:"-:migration;->"`
}

type OcUserActivity struct {
	ID        uint            `json:"-" gorm:"primaryKey;autoIncrement"`
	UID       string          `json:"uid" gorm:"type:varchar(26);not null;unique"`
	UserID    uint64          `json:"-" gorm:"index"`
	Log       json.RawMessage `json:"log" gorm:"type:json"`
	CreatedAt time.Time       `json:"createdAt" gorm:"autoCreateTime"`
}

type OcUserTrafficStatistics struct {
	ID     uint      `json:"-" gorm:"primaryKey;autoIncrement"`
	UID    string    `json:"uid" gorm:"type:varchar(26);not null;unique"`
	UserID uint64    `json:"-" gorm:"index"`
	Date   time.Time `json:"date" gorm:"date"`
	Rx     float64   `json:"rx" gorm:"numeric default 0.00"`
	Tx     float64   `json:"tx" gorm:"numeric default 0.00"`
}

func (o *OcUser) BeforeCreate(tx *gorm.DB) (err error) {
	if o.TrafficType == Free {
		o.TrafficSize = 0
	}
	o.UID = ulid.Make().String()
	return
}

func (o *OcUser) BeforeUpdate(tx *gorm.DB) (err error) {
	if o.TrafficType == Free {
		o.TrafficSize = 0
	}
	return
}

func (a *OcUserActivity) BeforeCreate(tx *gorm.DB) (err error) {
	a.UID = ulid.Make().String()
	return
}

func (s *OcUserTrafficStatistics) BeforeCreate(tx *gorm.DB) (err error) {
	s.UID = ulid.Make().String()
	return
}
