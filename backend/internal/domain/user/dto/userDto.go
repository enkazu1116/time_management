package dto

import (
	"time"

	User "time_management/backend/internal/domain/user/model"
	"time_management/backend/internal/infrastrcture/jst"
)

type CreateUserRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Password  string `json:"password"`
	Email     string `json:"email"`
}

type UserResponse struct {
	UserID    string    `json:"user_id"`
	UserName  string    `json:"user_name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewUserResponse(user *User.User) UserResponse {
	return UserResponse{
		UserID:    user.UserID.String(),
		UserName:  user.UserName.GetFullName(),
		Email:     user.Email.String(),
		CreatedAt: jst.In(user.CreatedAt),
		UpdatedAt: jst.In(user.UpdatedAt),
	}
}
