package ocUser

import (
	"api/internal/models"
	"api/pkg/utils"
)

type OcUsersResponse struct {
	OcUsers *[]models.OcUser
	Meta    *utils.ResponsePagination
}
