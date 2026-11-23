package di

import (
	"context"
	"fmt"

	userinterfaces "time_management/backend/domain/user/interfaces"
	userusecases "time_management/backend/domain/user/usecases"
	"time_management/backend/internal/db/connection"
	"time_management/backend/internal/db/sqlc"
	"time_management/backend/internal/repository"

	"github.com/jackc/pgx/v5/pgxpool"
)

/*
 * Container
 * アプリ全体の依存関係をここに集約する。
 * ドメインが増えたら Repositories / Usecases と NewContainer の組み立てだけ足す。
 */
type Container struct {
	DB           *pgxpool.Pool
	Repositories Repositories
	Usecases     Usecases
}

type Repositories struct {
	User userinterfaces.UserRepository
}

type Usecases struct {
	User *userusecases.UserUsecase
}

func NewContainer(ctx context.Context) (*Container, error) {
	pool, err := connection.DBConnect(ctx)
	if err != nil {
		return nil, fmt.Errorf("di: db connect: %w", err)
	}

	queries := sqlc.New(pool)

	userRepo := repository.NewUserRepository(queries)

	return &Container{
		DB: pool,
		Repositories: Repositories{
			User: userRepo,
		},
		Usecases: Usecases{
			User: userusecases.NewUserUsecase(userRepo),
		},
	}, nil
}

func (c *Container) Close() {
	if c == nil || c.DB == nil {
		return
	}
	c.DB.Close()
}
