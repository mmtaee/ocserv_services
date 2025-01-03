package ocGroup

import (
	"api/pkg/ocserv"
	"api/pkg/utils"
)

type GroupsResponse struct {
	Groups *[]ocserv.OcGroupConfig
	Meta   *utils.ResponsePagination
}
