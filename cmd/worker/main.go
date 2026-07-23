package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/jackc/pgx/v5/pgxpool"

	ownershipApp "github.com/Trivenqo/GhostKey/internal/ownership/application"
	ownershipInfra "github.com/Trivenqo/GhostKey/internal/ownership/infrastructure"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:5432/ghostkey?sslmode=disable"
	}

	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer pool.Close()

	kafkaBrokers := []string{os.Getenv("KAFKA_BROKER")}
	if len(kafkaBrokers) == 0 || kafkaBrokers[0] == "" {
		kafkaBrokers = []string{"localhost:9092"}
	}

	// Initialize Ownership Dependencies for Kafka Worker
	ownershipRepo := ownershipInfra.NewPostgresOwnershipRepository(pool)
	ownershipUseCase := ownershipApp.NewOwnershipUseCase(ownershipRepo)

	kafkaConsumer := ownershipInfra.NewKafkaConsumer(
		kafkaBrokers,
		"ghostkey.discovery.identity_discovered",
		"ghostkey.ownership.consumer-group",
		ownershipUseCase,
	)

	go kafkaConsumer.Start(ctx)

	// Graceful Shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Println("Shutting down worker...")
	cancel()
}