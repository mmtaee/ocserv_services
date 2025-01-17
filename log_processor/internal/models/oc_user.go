package models

import (
	"fmt"
	"gorm.io/gorm"
	"time"
)

const (
	Free            = "Free"
	MonthlyTransmit = "MonthlyTransmit"
	MonthlyReceive  = "MonthlyReceive"
	TotallyTransmit = "TotallyTransmit"
	TotallyReceive  = "TotallyReceive"
)

const (
	Connected    = "Connected"
	Disconnected = "Disconnected"
	Failed       = "Failed"
)

type OcUser struct {
	ID          uint       `json:"-" gorm:"primaryKey;autoIncrement"`
	UID         string     `json:"uid" gorm:"type:varchar(26);not null;unique"`
	Group       string     `json:"group" gorm:"type:varchar(16);default:'defaults'"`
	Username    string     `json:"username" gorm:"type:varchar(16);not null;unique"`
	Password    string     `json:"password" gorm:"type:varchar(16);not null"`
	IsLocked    bool       `json:"is_locked" gorm:"default(false)"`
	CreatedAt   time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
	ExpireAt    *time.Time `json:"expire_at"`
	TrafficType string     `json:"traffic_type" gorm:"type:varchar(32);not null;default:1" enums:"Free,MonthlyTransmit,MonthlyReceive,TotallyTransmit,TotallyReceive"`
	TrafficSize int        `json:"traffic_size" gorm:"not null;default:10"` // in GiB  >> x * 1024 ** 3
	Rx          int        `json:"rx" gorm:"not null;default:0"`            // Receive in bytes
	Tx          int        `json:"tx" gorm:"not null;default:0"`            // Transmit in bytes
	Description string     `json:"description" gorm:"type:text"`
	IsOnline    bool       `json:"is_online" gorm:"-:migration;->"`
}

type OcUserActivity struct {
	ID        uint      `json:"-" gorm:"primaryKey;autoIncrement"`
	OcUserID  uint      `json:"-" gorm:"index;constraint:OnDelete:CASCADE;"`
	Type      string    `json:"type" gorm:"type:varchar(32);not null;default:1" enums:"Connected,Disconnected,Failed"`
	Log       string    `json:"log" gorm:"type:text"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
}

type OcUserTrafficStatistics struct {
	ID        uint      `json:"-" gorm:"primaryKey;autoIncrement"`
	OcUserID  uint      `json:"-" gorm:"index;constraint:OnDelete:CASCADE"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	Rx        int       `json:"rx" gorm:"default:0"` // in bytes
	Tx        int       `json:"tx" gorm:"default:0"` // in bytes
}

func ValidateTrafficType(trafficType string) bool {
	switch trafficType {
	case Free, MonthlyTransmit, MonthlyReceive, TotallyTransmit, TotallyReceive:
		return true
	default:
		return false
	}
}

func (o *OcUser) BeforeSave(tx *gorm.DB) (err error) {
	if !ValidateTrafficType(o.TrafficType) {
		return fmt.Errorf("invalid TrafficType: %s", o.TrafficType)
	}
	if o.TrafficType == Free {
		o.TrafficSize = 0
	}
	return nil
}

func ValidateActivityType(t string) bool {
	switch t {
	case Connected, Disconnected, Failed:
		return true
	default:
		return false
	}
}

func (a *OcUserActivity) BeforeSave(tx *gorm.DB) (err error) {
	if !ValidateActivityType(a.Type) {
		return fmt.Errorf("invalid Type: %s", a.Type)
	}
	return nil
}
