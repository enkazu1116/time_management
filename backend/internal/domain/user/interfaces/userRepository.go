package interfaces

import (
	"context"
	User "time_management/backend/internal/domain/user/model"
	"uuid"
)

/*
 * User Repository
 * クエリとコマンド両方を記載する。
 * 分離したいタイミングでリファクタリングを行う。
 */
type UserRepository interface {
	GetUserById(ctx context.Context, id uuid.UUID) (*User.User, error)
	GetUserByUserName(ctx context.Context, userName string) ([]*User.User, error)
	GetUserByEmail(ctx context.Context, email string) (*User.User, error)
	GetUserByUserNameAndEmail(ctx context.Context, userName string, email string) (*User.User, error)
	CreateUser(ctx context.Context, user *User.User) error
	UpdateUser(ctx context.Context, user *User.User) error
	DeleteUser(ctx context.Context, id uuid.UUID) error
}
