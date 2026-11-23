package usecases

import (
	"context"
	"time_management/backend/domain/user/interfaces"
	User "time_management/backend/domain/user/model"
)

type UserUsecase struct {
	userRepository interfaces.UserRepository
}

/*
 * UserUsecase
 * ユーザー情報を管理するファクトリー関数
 * ファクトリー関数内で依存関係を注入し、以降の関数でロジックに集中
 */
func NewUserUsecase(r interfaces.UserRepository) *UserUsecase {
	return &UserUsecase{
		userRepository: r,
	}
}

func (u *UserUsecase) CreateUser(ctx context.Context, user *User.User) error {
	return u.userRepository.CreateUser(ctx, user)
}
