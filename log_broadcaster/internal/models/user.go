package models

import "time"

type User struct {
	ID        uint        `json:"-" gorm:"primaryKey;autoIncrement"`
	UID       string      `json:"uid" gorm:"type:varchar(26);not null;unique"`
	Username  string      `json:"username" gorm:"type:varchar(16);not null;unique"`
	Password  string      `json:"-" gorm:"type:varchar(64); not null"`
	IsAdmin   bool        `json:"is_admin" gorm:"type:bool;default(false)"`
	Salt      string      `json:"-" gorm:"type:varchar(6);"`
	LastLogin *time.Time  `json:"last_login"`
	Token     []UserToken `json:"-"`
}

type UserToken struct {
	ID        uint       `json:"-" gorm:"primaryKey;autoIncrement"`
	UserID    uint       `json:"-" gorm:"index"`
	UID       string     `json:"uid" gorm:"type:varchar(26);not null;unique"`
	Token     string     `json:"token" gorm:"type:varchar(128)"`
	CreatedAt time.Time  `json:"created_at" gorm:"autoCreateTime"`
	ExpireAt  *time.Time `json:"expire_at"`
	User      User       `json:"user"`
}
