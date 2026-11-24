package connection

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

const (
	DB_CONNECTION_ERROR = "DB接続に必要な情報が不足しています"
)

const (
	DB_USER     = "DB_USER"
	DB_PASSWORD = "DB_PASSWORD"
	DB_HOST     = "DB_HOST"
	DB_PORT     = "DB_PORT"
	DB_NAME     = "DB_NAME"
)

func DBConnect(ctx context.Context) (*pgxpool.Pool, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("%s: %w", DB_CONNECTION_ERROR, err)
	}

	dbUser := os.Getenv(DB_USER)
	dbPassword := os.Getenv(DB_PASSWORD)
	dbHost := os.Getenv(DB_HOST)
	dbPort := os.Getenv(DB_PORT)
	dbName := os.Getenv(DB_NAME)

	if dbUser == "" || dbPassword == "" || dbHost == "" || dbPort == "" || dbName == "" {
		return nil, fmt.Errorf("%s", DB_CONNECTION_ERROR)
	}

	dsn := fmt.Sprintf(
		"user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
		dbUser, dbPassword, dbHost, dbPort, dbName,
	)

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", DB_CONNECTION_ERROR, err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("%s: %w", DB_CONNECTION_ERROR, err)
	}

	return pool, nil
}
