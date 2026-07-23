package sdk

import (
	"context"
	"time"

	"go.uber.org/zap"
)

// RecordHandler is a callback provided by the host module (e.g., Discovery)
// to process the raw records fetched by the scheduler.
type RecordHandler func(ctx context.Context, connectorName string, records []RawRecord) error

// Scheduler orchestrates the polling loop for all registered connectors.
type Scheduler struct {
	registry     *Registry
	checkpointer CheckpointManager
	credManager  CredentialManager
	logger       *zap.Logger
}

// NewScheduler creates a new scheduler instance.
func NewScheduler(
	registry *Registry,
	checkpointer CheckpointManager,
	credManager CredentialManager,
	logger *zap.Logger,
) *Scheduler {
	return &Scheduler{
		registry:     registry,
		checkpointer: checkpointer,
		credManager:  credManager,
		logger:       logger,
	}
}

// StartWorker begins a polling loop for a specific connector.
func (s *Scheduler) StartWorker(ctx context.Context, connectorName string, interval time.Duration, handler RecordHandler) {
	connector, err := s.registry.Get(connectorName)
	if err != nil {
		s.logger.Error("Failed to start worker: connector not found", zap.String("connector", connectorName))
		return
	}

	s.logger.Info("Starting connector worker", zap.String("connector", connectorName), zap.Duration("interval", interval))

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	// Trigger an immediate initial sync before waiting for the first tick
	s.runSync(ctx, connector, handler)

	for {
		select {
		case <-ctx.Done():
			s.logger.Info("Stopping connector worker", zap.String("connector", connectorName))
			return
		case <-ticker.C:
			s.runSync(ctx, connector, handler)
		}
	}
}

func (s *Scheduler) runSync(ctx context.Context, connector Connector, handler RecordHandler) {
	name := connector.Name()
	log := s.logger.With(zap.String("connector", name))

	log.Debug("Starting sync run")

	// 1. Fetch Credentials
	creds, err := s.credManager.Get(ctx, name)
	if err != nil {
		log.Error("Failed to retrieve credentials", zap.Error(err))
		return
	}

	// 2. Authenticate
	if err := connector.Authenticate(ctx, creds); err != nil {
		log.Error("Failed to authenticate connector", zap.Error(err))
		return
	}

	// 3. Load Checkpoint
	cursor, err := s.checkpointer.Load(ctx, name)
	if err != nil {
		log.Error("Failed to load checkpoint", zap.Error(err))
		return
	}

	// 4. Paginate and Discover
	for {
		if ctx.Err() != nil {
			log.Warn("Sync interrupted by context cancellation")
			return
		}

		page, err := connector.Discover(ctx, cursor)
		if err != nil {
			log.Error("Discovery failed during pagination", zap.Error(err), zap.String("token", cursor.Token))
			return
		}

		if len(page.Items) > 0 {
			// 5. Pass records to the host module for normalization
			if err := handler(ctx, name, page.Items); err != nil {
				log.Error("Handler failed to process records", zap.Error(err))
				return
			}
		}

		// 6. Save Checkpoint safely
		if err := s.checkpointer.Save(ctx, name, page.NextCursor); err != nil {
			log.Error("Failed to save checkpoint", zap.Error(err))
			return
		}

		if !page.HasMore {
			break
		}
		
		cursor = page.NextCursor
	}

	log.Debug("Sync run completed successfully")
}