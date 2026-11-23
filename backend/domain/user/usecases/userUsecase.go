package usecases

import (
	"context"
	"fmt"
	"time"
	"uuid"

	"time_management/backend/domain/user/dto"
	"time_management/backend/domain/user/interfaces"
	User "time_management/backend/domain/user/model"
	valueobjects "time_management/backend/domain/user/model/valueObjects"
	"time_management/backend/internal/util/jst"
)

type UserUsecase struct {
	userRepository interfaces.UserRepository
}

func NewUserUsecase(r interfaces.UserRepository) *UserUsecase {
	return &UserUsecase{
		userRepository: r,
	}
}

func (u *UserUsecase) CreateUser(ctx context.Context, req dto.CreateUserRequest) (*dto.UserResponse, error) {
	now := jst.Now()
	user := User.NewUser(
		uuid.New(),
		valueobjects.NewUserName(req.FirstName, req.LastName),
		valueobjects.NewPassword(req.Password),
		valueobjects.NewEmailAddress(req.Email),
		now,
		now,
		time.Time{},
	)

	if err := u.userRepository.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	res := dto.NewUserResponse(user)
	return &res, nil
}

func (u *UserUsecase) GetUserByID(ctx context.Context, id string) (*dto.UserResponse, error) {
	userID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid user id: %w", err)
	}

	user, err := u.userRepository.GetUserById(ctx, userID)
	if err != nil {
		return nil, err
	}

	res := dto.NewUserResponse(user)
	return &res, nil
}
