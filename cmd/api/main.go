package main

import (
	"context"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"

	ownershipApp "github.com/Trivenqo/GhostKey/internal/ownership/application"
	ownershipInfra "github.com/Trivenqo/GhostKey/internal/ownership/infrastructure"
	ownershipPres "github.com/Trivenqo/GhostKey/internal/ownership/presentation"
)

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:5432/ghostkey?sslmode=disable"
	}

	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer pool.Close()

	app := fiber.New()
	v1 := app.Group("/v1")

	// Ownership Module Initialization
	ownershipRepo := ownershipInfra.NewPostgresOwnershipRepository(pool)
	ownershipUseCase := ownershipApp.NewOwnershipUseCase(ownershipRepo)
	ownershipHandler := ownershipPres.NewOwnershipHandler(ownershipUseCase)
	ownershipHandler.RegisterRoutes(v1)

	log.Println("API server listening on :8080")
	if err := app.Listen(":8080"); err != nil {
		log.Fatalf("Fiber error: %v", err)
	}
}