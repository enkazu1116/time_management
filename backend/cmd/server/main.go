package main

import (
	"context"
	"log"
	"net/http"

	"time_management/backend/internal/infrastrcture/di"
)

func main() {
	ctx := context.Background()

	container, err := di.NewContainer(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer container.Close()

	if err := http.ListenAndServe(":8080", container.Router); err != nil {
		log.Fatal(err)
	}

}
