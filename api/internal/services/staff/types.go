package staff

import (
	"api/internal/models"
	"api/pkg/utils"
)

type StaffsResponse struct {
	Staffs *[]models.User
	Meta   *utils.ResponsePagination
}

type CreateStaffRequest struct {
	User struct {
		Username string `json:"username" validate:"required,min=2,max=16" example:"john_doe" `
		Password string `json:"password" validate:"required,min=2,max=16" example:"doe123456"`
	} `json:"user" validator:"required"`
	Permission struct {
		OcUser    bool `json:"oc_user" validator:"required"`
		OcGroup   bool `json:"oc_group" validator:"required"`
		Statistic bool `json:"statistic" validator:"required"`
		Occtl     bool `json:"occtl" validator:"required"`
		System    bool `json:"system" validator:"required"`
	} `json:"permission" validator:"required"`
}

type UpdateStaffPermissionRequest struct {
	OcUser    bool `json:"oc_user" validator:"required"`
	OcGroup   bool `json:"oc_group" validator:"required"`
	Statistic bool `json:"statistic" validator:"required"`
	Occtl     bool `json:"occtl" validator:"required"`
	System    bool `json:"system" validator:"required"`
}

type UpdateStaffPasswordRequest struct {
	Password string `json:"password" validate:"required,min=2,max=16" example:"doe123456"`
}