package models

import "time"

type User struct {
	ID         uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	Username   string         `json:"username" gorm:"type:varchar(16);not null;unique"`
	Password   string         `json:"password" gorm:"type:varchar(16) not null"`
	IsAdmin    bool           `json:"is_admin" gorm:"type:bool;default(false)"`
	LastLogin  time.Time      `json:"last_login"`
	Token      []UserToken    `json:"tokens"`
	Permission UserPermission `json:"permission"`
}

type UserPermission struct {
	ID        uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID    uint64 `json:"_" gorm:"index"`
	OcUser    bool   `json:"oc_user" gorm:"default(false)"`
	OcGroup   bool   `json:"oc_group" gorm:"default(false)"`
	Statistic bool   `json:"statistic" gorm:"default(false)"`
	Occtl     bool   `json:"occtl" gorm:"default(false)"`
	System    bool   `json:"system" gorm:"default(false)"`
}

type UserToken struct {
	ID        uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID    uint64    `json:"_" gorm:"index"`
	Token     string    `json:"token" gorm:"type:varchar(128)"`
	CreatedAt time.Time `json:"createdAt" gorm:"autoCreateTime"`
	ExpiresAt time.Time `json:"expiresAt"`
	User      User      `json:"user"`
}
