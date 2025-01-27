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
// - update_oc_default_group
// - create_oc_group
// - update_oc_group
// - delete_oc_group

// *oc user
// - create_oc_user

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
