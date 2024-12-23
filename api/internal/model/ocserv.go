package model

import "time"

const (
	Free int32 = iota
	Monthly
	Totally
)

type User struct {
	ID          uint64    `json:"id" xorm:"primary_key;auto_increment"`
	Group       string    `json:"group" xorm:"varchar(16)" default:"defaults"`
	Username    string    `json:"username" xorm:"varchar(16) not null unique"`
	Password    string    `json:"password" xorm:"varchar(16) not null"`
	IsLocked    bool      `json:"isLocked" xorm:"default false"`
	CreatedAt   time.Time `json:"createdAt" xorm:"created"`
	UpdatedAt   time.Time `json:"updatedAt" xorm:"updated"`
	ExpiresAt   time.Time `json:"expiresAt"`
	TrafficType int32     `json:"trafficType" xorm:"tinyint default 0"`
	TrafficSize int32     `json:"trafficSize" xorm:"tinyint default 10"`
	Rx          float64   `json:"rx" xorm:"numeric default 0.00"`
	Tx          float64   `json:"tx" xorm:"numeric default 0.00"`
	Description string    `json:"description" xorm:"text"`
}

func (u *User) BeforeInsert() {
	if u.TrafficType == Free {
		u.TrafficSize = 0
	}
}

func (u *User) BeforeUpdate() {
	if u.TrafficType == Free {
		u.TrafficSize = 0
	}
}

type TrafficStatistics struct {
	UserID uint64    `json:"userId" xorm:"index"`
	Date   time.Time `json:"date" xorm:"date"`
	Rx     float64   `json:"rx" xorm:"numeric default 0.00"`
	Tx     float64   `json:"tx" xorm:"numeric default 0.00"`
}
