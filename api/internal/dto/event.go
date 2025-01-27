package dto

import "github.com/mmtaee/go-oc-utils/models"

// events :

// *staff :
//	- create_staff
//	- create_staff_permission
//	- update_staff_permission
//	- update_staff_password
//	- delete_staff

// *panel:
//	- update_panel_config

// *oc group
// - update_default_group
// - create_group
// - update_group
// - delete_group

type CreateStaffEvent struct {
	User models.User `json:"user"`
}

type CreatePermissionEvent struct {
	Permission models.UserPermission `json:"permission"`
}

type UpdateStaffPermissionEvent struct {
	Permission models.UserPermission `json:"permission"`
}

type UpdatePanelConfigEvent struct {
	Config models.PanelConfig `json:"config"`
}
