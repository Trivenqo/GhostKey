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
	"go.uber.org/zap"
)

// mockCredManager simulates fetching secrets from Vault/AWS Secrets Manager
type mockCredManager struct{}

func (m *mockCredManager) Get(ctx context.Context, connectorName string) (sdk.Credentials, error) {
	// Provide the fake credentials the AWS Connector expects to pass authentication
	return sdk.Credentials{
		"access_key_id":     "AKIAIOSFODNN7EXAMPLE",
		"secret_access_key": "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
	}, nil
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 1. Boot infrastructure (Redis is needed for Checkpoints)
	container, err := bootstrap.BuildContainer(ctx)
	if err != nil {
		panic(err)
	}
	defer container.Close()

	logger := container.Logger
	logger.Info("Starting GhostKey Background Worker")

	// 2. Setup the Connector SDK
	registry := sdk.NewRegistry()
	if err := registry.Register(aws.NewConnector()); err != nil {
		logger.Fatal("Failed to register AWS connector", zap.Error(err))
	}

	checkpointManager := sdk.NewRedisCheckpointManager(container.Redis)
	credManager := &mockCredManager{}

	scheduler := sdk.NewScheduler(registry, checkpointManager, credManager, logger)

	// 3. Define the RecordHandler
	// In the real app, this handler will publish events to Kafka.
	// For now, we just normalize and print them!
	awsNormalizer := aws.NewNormalizer()

	handler := func(ctx context.Context, connectorName string, records []sdk.RawRecord) error {
		for _, rec := range records {
			// Pass raw JSON through the Normalizer
			identity, err := awsNormalizer.Normalize(rec)
			if err != nil {
				logger.Error("Failed to normalize record", zap.Error(err))
				continue
			}

			// We now have a clean canonical.Identity!
			logger.Info("✅ Discovered Canonical Identity",
				zap.String("connector", connectorName),
				zap.String("type", string(identity.Type)),
				zap.String("name", identity.DisplayName),
				zap.String("external_id", identity.ExternalRef.ExternalID),
			)
		}
		return nil
	}

	// 4. Start the Polling Loop in the background (runs every 10 seconds for testing)
	go scheduler.StartWorker(ctx, "aws", 10*time.Second, handler)

	// 5. Wait for termination signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down worker...")
	cancel() // Stops the scheduler loops

	// Give background tasks a moment to finish current page
	time.Sleep(1 * time.Second)
	logger.Info("Worker stopped")
}