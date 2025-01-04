package ocUser

import (
	"api/internal/models"
	"api/pkg/utils"
	"time"
)

type OcservUsersResponse struct {
	OcUsers *[]models.OcUser
	Meta    *utils.ResponsePagination
}

type OcservUserCreateRequest struct {
	Group       string     `json:"group" validate:"required"`
	Username    string     `json:"username" validate:"required;min=3,max=16"`
	Password    string     `json:"password" validate:"required,min=8,max=16"`
	TrafficType int32      `json:"trafficType" validate:"required;min=0,max=2"`
	TrafficSize int32      `json:"trafficSize"`
	ExpireAt    *time.Time `json:"expire_at"`
}
