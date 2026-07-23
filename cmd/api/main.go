package main

import (
	"context"
	"log"
	"os"

	"github.com/Trivenqo/GhostKey/internal/bootstrap"
	"github.com/Trivenqo/GhostKey/internal/interfaces/http"
	"go.uber.org/zap"
)

func main() {
	ctx := context.Background()

	// 1. Build Infrastructure Container
	container, err := bootstrap.BuildContainer(ctx)
	if err != nil {
		log.Fatalf("Failed to initialize container: %v", err)
	}
	defer container.Close()

	// 2. Build HTTP Router
	router := http.NewRouter(container)

	// 3. Assemble and Start App
	app := bootstrap.NewApp(container, router)
	
	if err := app.Start(); err != nil {
		container.Logger.Fatal("Application terminated", zap.Error(err))
		os.Exit(1)
	}
}