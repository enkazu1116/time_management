package model

import (
	"time"
	model "time_management/backend/domain/shared"
	valueobjects "time_management/backend/domain/user/model/valueObjects"
)

/*
 * User Domain Model
 * AuditLogはアーキテクチャ上必要となるため、ドメインに含める
 */
type User struct {
	UserID    model.UUID                 `json:"user_id"`
	UserName  *valueobjects.UserName     `json:"user_name"`
	Password  *valueobjects.Password     `json:"password"`
	Email     *valueobjects.EmailAddress `json:"email"`
	CreatedAt time.Time                  `json:"created_at"`
	UpdatedAt time.Time                  `json:"updated_at"`
	DeletedAt time.Time                  `json:"deleted_at"`
}

func NewUser(
	userID model.UUID,
	userName *valueobjects.UserName,
	password *valueobjects.Password,
	email *valueobjects.EmailAddress,
	createdAt time.Time,
	updatedAt time.Time,
	deletedAt time.Time,
) *User {
	return &User{
		UserID:    userID,
		UserName:  userName,
		Password:  password,
		Email:     email,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
		DeletedAt: deletedAt,
	}
}
