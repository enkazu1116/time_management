package connection

import (
	"database/sql"
	"fmt"
	"log"
	"os"

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

type DBServer struct {
	DB *sql.DB
}

func DBConnect() (*sql.DB, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(DB_CONNECTION_ERROR)
		return nil, err
	}

	dbUser := os.Getenv(DB_USER)
	dbPassword := os.Getenv(DB_PASSWORD)
	dbHost := os.Getenv(DB_HOST)
	dbPort := os.Getenv(DB_PORT)
	dbName := os.Getenv(DB_NAME)

	db, err := sql.Open("postgres", fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable", dbUser, dbPassword, dbHost, dbPort, dbName))
	if err != nil {
		log.Fatal(DB_CONNECTION_ERROR)
		return nil, err
	}

	return db, nil
}
