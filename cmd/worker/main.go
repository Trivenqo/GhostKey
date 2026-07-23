package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Trivenqo/GhostKey/internal/bootstrap"
	"github.com/Trivenqo/GhostKey/internal/connector/aws"
	"github.com/Trivenqo/GhostKey/internal/connector/sdk"
	
	"github.com/Trivenqo/GhostKey/internal/discovery/application/usecase"
	"github.com/Trivenqo/GhostKey/internal/discovery/infrastructure/kafka"
	"github.com/Trivenqo/GhostKey/internal/discovery/infrastructure/postgres"

	"go.uber.org/zap"
)

type mockCredManager struct{}

func (m *mockCredManager) Get(ctx context.Context, connectorName string) (sdk.Credentials, error) {
	return sdk.Credentials{
		"access_key_id":     "AKIAIOSFODNN7EXAMPLE",
		"secret_access_key": "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
	}, nil
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	container, err := bootstrap.BuildContainer(ctx)
	if err != nil {
		panic(err)
	}
	defer container.Close()

	logger := container.Logger
	logger.Info("Starting GhostKey Background Worker")

	// 1. Setup Discovery Module Dependencies
	repo := postgres.NewIdentityRepository(container.DB)
	pub := kafka.NewPublisher(container.Kafka)
	registerUseCase := usecase.NewRegisterIdentityUseCase(repo, pub)

	// 2. Setup the Connector SDK
	registry := sdk.NewRegistry()
	if err := registry.Register(aws.NewConnector()); err != nil {
		logger.Fatal("Failed to register AWS connector", zap.Error(err))
	}

	checkpointManager := sdk.NewRedisCheckpointManager(container.Redis)
	credManager := &mockCredManager{}
	scheduler := sdk.NewScheduler(registry, checkpointManager, credManager, logger)

	// 3. Define the RecordHandler using the Use Case
	awsNormalizer := aws.NewNormalizer()

	handler := func(ctx context.Context, connectorName string, records []sdk.RawRecord) error {
		for _, rec := range records {
			identity, err := awsNormalizer.Normalize(rec)
			if err != nil {
				logger.Error("Failed to normalize record", zap.Error(err))
				continue
			}

			// Hand off to the Business Logic!
			if err := registerUseCase.Execute(ctx, identity); err != nil {
				logger.Error("Failed to register identity", zap.Error(err))
				// We return the error so the Checkpoint is NOT saved if DB/Kafka is down.
				return err
			}
			
			logger.Info("✅ Discovered and Persisted Identity",
				zap.String("name", identity.DisplayName),
				zap.String("external_id", identity.ExternalRef.ExternalID),
			)
		}
		return nil
	}

	// 4. Start the Polling Loop
	go scheduler.StartWorker(ctx, "aws", 10*time.Second, handler)

	// 5. Wait for termination signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down worker...")
	cancel()
	time.Sleep(1 * time.Second)
}