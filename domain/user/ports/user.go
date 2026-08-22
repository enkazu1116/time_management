package ports

import (
	"errors"

	"time_management/domain/shared/messages"
	sharedmodel "time_management/domain/shared/model"
)

var (
	ErrUserIDRequired   = errors.New(messages.UserIDRequired)
	ErrUserNameRequired = errors.New(messages.UserNameRequired)
	ErrRoleRequired     = errors.New(messages.RoleRequired)
)

type RegisterUserInput struct {
	UserID       string
	Password     string
	Name         string
	Email        string
	RegularStart string
	RegularEnd   string
	RoleID       string
	RoleName     string
}

func (i RegisterUserInput) Validate() error {
	if i.UserID == "" {
		return ErrUserIDRequired
	}
	if err := sharedmodel.ValidateUUIDString(i.UserID); err != nil {
		return err
	}
	if i.Name == "" {
		return ErrUserNameRequired
	}
	if i.RoleID == "" || i.RoleName == "" {
		return ErrRoleRequired
	}
	if err := sharedmodel.ValidateUUIDString(i.RoleID); err != nil {
		return err
	}
	return nil
}
