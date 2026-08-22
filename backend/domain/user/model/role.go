package model

import sharedmodel "time_management/domain/shared/model"

type Role struct {
	RoleID   string
	RoleName string
}

func NewRole(roleID string, roleName string) Role {
	return Role{RoleID: roleID, RoleName: roleName}
}

func (r Role) Validate() error {
	if r.RoleID == "" || r.RoleName == "" {
		return ErrRoleRequired
	}
	if err := sharedmodel.ValidateUUIDString(r.RoleID); err != nil {
		return err
	}
	return nil
}
