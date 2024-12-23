package model

type Admin struct {
	ID          uint   `json:"id" xorm:"primary_key;auto_increment"`
	Username    string `json:"username" xorm:"varchar(16) not null unique"`
	Password    string `json:"password" xorm:"varchar(16) not null"`
	IsSuperuser bool   `json:"is_superuser" xorm:"default false"`
	LastLogin   string `json:"last_login" xorm:"datetime default CURRENT_TIMESTAMP"`
}

type Permission struct {
	ID         uint   `json:"id" xorm:"primary_key;auto_increment"`
	AdminID    uint64 `json:"adminId" xorm:"index"`
	Users      bool   `json:"users" xorm:"default false"`
	Groups     bool   `json:"groups" xorm:"default false"`
	Statistics bool   `json:"statistics" xorm:"default false"`
	Occtl      bool   `json:"occtl" xorm:"default false"`
}
