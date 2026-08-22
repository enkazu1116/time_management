package repositories

import "time_management/domain/user/model"

type UserRepository interface {
	Create(user model.User) error
	Update(user model.User) error
	Delete(userID string) error
	FindByUserID(userID string) (model.User, error)
}
