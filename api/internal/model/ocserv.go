package model

import (
	"gorm.io/gorm"
	"time"
)

const (
	Free int32 = iota
	Monthly
	Totally
)

type User struct {
	ID          uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Group       string    `json:"group" gorm:"type:varchar(16);default('defaults')"`
	Username    string    `json:"username" gorm:"type:varchar(16);not null;unique"`
	Password    string    `json:"password" gorm:"type:varchar(16);not null"`
	IsLocked    bool      `json:"isLocked" gorm:"default(false)"`
	CreatedAt   time.Time `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updatedAt" gorm:"autoUpdateTime"`
	ExpiresAt   time.Time `json:"expiresAt"`
	TrafficType int32     `json:"trafficType" gorm:"not null;default(0)"`
	TrafficSize int32     `json:"trafficSize" gorm:"not null;default(10)"`
	Rx          float64   `json:"rx" gorm:"not null;default(0.00)"`
	Tx          float64   `json:"tx" gorm:"not null;default(0.00)"`
	Description string    `json:"description" gorm:"type:text"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.TrafficType == Free {
		u.TrafficSize = 0
	}
	return
}

func (u *User) BeforeUpdate(tx *gorm.DB) (err error) {
	if u.TrafficType == Free {
		u.TrafficSize = 0
	}
	return
}

type TrafficStatistics struct {
	UserID uint64    `json:"userId" gorm:"index"`
	Date   time.Time `json:"date" gorm:"date"`
	Rx     float64   `json:"rx" gorm:"numeric default 0.00"`
	Tx     float64   `json:"tx" gorm:"numeric default 0.00"`
}

type GroupConfig struct {
	RxDataPerSec         string   `json:"rx-data-per-sec"`
	TxDataPerSec         string   `json:"tx-data-per-sec"`
	MaxSameClients       int      `json:"max-same-clients"`
	IPv4Network          string   `json:"ipv4-network"`
	DNS                  []string `json:"dns"`
	NoUDP                bool     `json:"no-udp"`
	KeepAlive            int      `json:"keepalive"`
	DPD                  int      `json:"dpd"`
	MobileDPD            int      `json:"mobile-dpd"`
	TunnelAllDNS         bool     `json:"tunnel-all-dns"`
	RestrictUserToRoutes bool     `json:"restrict-user-to-routes"`
	StatsReportTime      int      `json:"stats-report-time"`
	MTU                  int      `json:"mtu"`
	IdleTimeout          int      `json:"idle-timeout"`
	MobileIdleTimeout    int      `json:"mobile-idle-timeout"`
	SessionTimeout       int      `json:"session-timeout"`
	//NoRoutes             bool     `json:"no_routes"`
	//Routes               []string `json:"routes"`
}
