package ocUser

import (
	"api/internal/models"
	"api/pkg/utils"
)

type OcservUsersResponse struct {
	OcUsers *[]models.OcUser          `json:"oc_users"`
	Meta    *utils.ResponsePagination `json:"meta"`
}

type OcservUserCreateOrUpdateRequest struct {
	Group       *string `json:"group" validate:"required"`
	Username    *string `json:"username" validate:"required,min=3,max=16"`
	Password    *string `json:"password" validate:"required,min=1,max=16"`
	TrafficType *string `json:"traffic_type" validate:"required" enums:"Free,MonthlyTransmit,MonthlyReceive,TotallyTransmit,TotallyReceive"`
	TrafficSize *int    `json:"traffic_size"`
	ExpireAt    *string `json:"expire_at" validate:"required"`
}

type OcservUserLockRequest struct {
	Lock *bool `json:"lock" validate:"required"`
}
