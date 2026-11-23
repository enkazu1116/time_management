package di

import (
	"context"
	"fmt"
	"net/http"

	userinterfaces "time_management/backend/domain/user/interfaces"
	userusecases "time_management/backend/domain/user/usecases"
	"time_management/backend/internal/db/connection"
	"time_management/backend/internal/db/sqlc"
	apihttp "time_management/backend/internal/http"
	"time_management/backend/internal/http/handlers"

	"github.com/jackc/pgx/v5/pgxpool"
)

/*
 * Container
 * アプリ全体の依存関係をここに集約する。
 * ドメインが増えたら ○○Infrastructure を追加して NewContainer に1行足す。
 */
type Container struct {
	DB           *pgxpool.Pool
	Repositories Repositories
	Usecases     Usecases
	Handlers     Handlers
	Router       http.Handler
}

type Repositories struct {
	User userinterfaces.UserRepository
}

type Usecases struct {
	User *userusecases.UserUsecase
}

type Handlers struct {
	User *handlers.UserHandler
}

func NewContainer(ctx context.Context) (*Container, error) {
	pool, err := connection.DBConnect(ctx)
	if err != nil {
		return nil, fmt.Errorf("di: db connect: %w", err)
	}

	queries := sqlc.New(pool)
	userInfra := NewUsersInfrastructure(queries)

	return &Container{
		DB: pool,
		Repositories: Repositories{
			User: userInfra.Repo,
		},
		Usecases: Usecases{
			User: userInfra.Usecase,
		},
		Handlers: Handlers{
			User: userInfra.Handler,
		},
		Router: apihttp.NewRouter(userInfra.Handler),
	}, nil
}

func (c *Container) Close() {
	if c == nil || c.DB == nil {
		return
	}
	c.DB.Close()
}
