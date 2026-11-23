package main

import (
	"log"
	"time_management/backend/internal/db/connection"
)

func main() {

	db, err := connection.DBConnect()
	if err != nil {
		log.Fatal(err)
		return
	}
	defer db.Close()

}
