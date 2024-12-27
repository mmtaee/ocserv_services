package models

import "time"

type Admin struct {
	ID          uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Username    string    `json:"username" gorm:"type:varchar(16);not null;unique"`
	Password    string    `json:"password" gorm:"type:varchar(16) not null"`
	IsSuperuser bool      `json:"is_superuser" gorm:"type:bool;default(false)"`
	LastLogin   time.Time `json:"last_login"`
}

type Permission struct {
	ID         uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	AdminID    uint64 `json:"adminId" gorm:"index"`
	Users      bool   `json:"users" gorm:"default(false)"`
	Groups     bool   `json:"groups" gorm:"default(false)"`
	Statistics bool   `json:"statistics" gorm:"default(false)"`
	Occtl      bool   `json:"occtl" gorm:"default(false)"`
}
