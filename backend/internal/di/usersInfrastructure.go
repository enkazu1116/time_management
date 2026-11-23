package di

import (
	"time_management/backend/domain/user/interfaces"
	"time_management/backend/domain/user/usecases"
	"time_management/backend/internal/db/sqlc"
	"time_management/backend/internal/http/handlers"
	"time_management/backend/internal/repository"
)

/*
 * UsersInfrastructure
 * User ドメインの repo → usecase → handler をまとめて組み立てる。
 */
type UsersInfrastructure struct {
	Repo    interfaces.UserRepository
	Usecase *usecases.UserUsecase
	Handler *handlers.UserHandler
}

func NewUsersInfrastructure(queries *sqlc.Queries) *UsersInfrastructure {
	repo := repository.NewUserRepository(queries)
	usecase := usecases.NewUserUsecase(repo)
	handler := handlers.NewUserHandler(usecase)

	return &UsersInfrastructure{
		Repo:    repo,
		Usecase: usecase,
		Handler: handler,
	}
}
