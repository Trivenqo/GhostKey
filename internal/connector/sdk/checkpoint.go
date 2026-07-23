package sdk

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"
)

// CheckpointManager handles saving and loading pagination state.
type CheckpointManager interface {
	Load(ctx context.Context, connector string) (Cursor, error)
	Save(ctx context.Context, connector string, cursor Cursor) error
}

// RedisCheckpointManager implements CheckpointManager using Redis.
type RedisCheckpointManager struct {
	client *redis.Client
	prefix string
}

func NewRedisCheckpointManager(client *redis.Client) *RedisCheckpointManager {
	return &RedisCheckpointManager{
		client: client,
		prefix: "ghostkey:checkpoint:", // e.g., ghostkey:checkpoint:aws
	}
}

func (m *RedisCheckpointManager) Load(ctx context.Context, connector string) (Cursor, error) {
	key := m.prefix + connector
	
	val, err := m.client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		// No checkpoint found, return an empty cursor (start from beginning)
		return Cursor{}, nil
	} else if err != nil {
		return Cursor{}, fmt.Errorf("failed to load checkpoint from redis: %w", err)
	}

	var cursor Cursor
	if err := json.Unmarshal([]byte(val), &cursor); err != nil {
		return Cursor{}, fmt.Errorf("failed to parse checkpoint data: %w", err)
	}

	return cursor, nil
}

func (m *RedisCheckpointManager) Save(ctx context.Context, connector string, cursor Cursor) error {
	key := m.prefix + connector

	data, err := json.Marshal(cursor)
	if err != nil {
		return fmt.Errorf("failed to serialize cursor: %w", err)
	}

	// Save with no expiration; cursors should persist between worker restarts
	if err := m.client.Set(ctx, key, data, 0).Err(); err != nil {
		return fmt.Errorf("failed to save checkpoint to redis: %w", err)
	}

	return nil
}