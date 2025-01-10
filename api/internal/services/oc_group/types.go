package ocGroup

import (
	"api/pkg/ocserv"
	"api/pkg/utils"
)

type GroupsResponse struct {
	Groups *[]ocserv.OcGroupConfig   `json:"groups"`
	Meta   *utils.ResponsePagination `json:"meta"`
}

type CreateGroupRequest struct {
	Name   string                `json:"name" validate:"required"`
	Config *ocserv.OcGroupConfig `json:"config" validate:"required"`
}
