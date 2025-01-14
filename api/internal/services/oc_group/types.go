package ocGroup

import (
	"api/pkg/utils"
	"github.com/mmtaee/go-oc-utils/handler/ocgroup"
)

type GroupsResponse struct {
	Groups *[]ocgroup.OcservGroupConfig `json:"groups"`
	Meta   *utils.ResponsePagination    `json:"meta"`
}

type CreateGroupRequest struct {
	Name   string                     `json:"name" validate:"required"`
	Config *ocgroup.OcservGroupConfig `json:"config" validate:"required"`
}
