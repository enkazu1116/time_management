package requests

import (
	"errors"

	sharedmodel "time_management/domain/shared/model"
	"time_management/infrastructure/util/messages"
)

var (
	ErrNameRequired     = errors.New(messages.RequestNameRequired)
	ErrPasswordRequired = errors.New(messages.RequestPasswordRequired)
	ErrEmailRequired    = errors.New(messages.RequestEmailRequired)
)

type RegisterUserRequest struct {
	UserID       string `json:"user_id"`
	Password     string `json:"password"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	RegularStart string `json:"regular_start"`
	RegularEnd   string `json:"regular_end"`
	RoleID       string `json:"role_id"`
	RoleName     string `json:"role_name"`
}

func (r RegisterUserRequest) Validate() error {
	if r.UserID != "" {
		if err := sharedmodel.ValidateUUIDString(r.UserID); err != nil {
			return err
		}
	}
	if r.Name == "" {
		return ErrNameRequired
	}
	if r.Password == "" {
		return ErrPasswordRequired
	}
	if r.Email == "" {
		return ErrEmailRequired
	}
	if r.RoleID != "" {
		if err := sharedmodel.ValidateUUIDString(r.RoleID); err != nil {
			return err
		}
	}
	return nil
}
