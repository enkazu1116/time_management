package usecases

import (
	"time_management/domain/user/interfaces/repositories"
	"time_management/domain/user/model"
	"time_management/domain/user/ports"
)

type UserUsecase interface {
	RegisterUser(input ports.RegisterUserInput) (model.User, error)
}

type userUsecase struct {
	userRepository repositories.UserRepository
}

func NewUserUsecase(userRepository repositories.UserRepository) UserUsecase {
	return &userUsecase{userRepository: userRepository}
}

func (u *userUsecase) RegisterUser(input ports.RegisterUserInput) (model.User, error) {
	if err := input.Validate(); err != nil {
		return model.User{}, err
	}

	user := model.NewUserFromParams(model.UserParams{
		UserID:       input.UserID,
		Password:     input.Password,
		Name:         input.Name,
		Email:        input.Email,
		RegularStart: input.RegularStart,
		RegularEnd:   input.RegularEnd,
		Role:         model.NewRole(input.RoleID, input.RoleName),
	})

	// 登録前にモデル自身へ妥当性確認を委譲し、利用者が個別ルールを知らなくてよい形にする。
	if err := user.Validate(); err != nil {
		return model.User{}, err
	}

	// 永続化の技術詳細は repository 側へ隔離し、ユースケースは登録手順だけを持つ。
	if err := u.userRepository.Create(user); err != nil {
		return model.User{}, err
	}

	return user, nil
}
