package main

import (
	"context"
	"log"

	"time_management/backend/internal/di"
)

func main() {
	ctx := context.Background()

	container, err := di.NewContainer(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer container.Close()

	// 以降は container から依存を取り出す
	_ = container
}
