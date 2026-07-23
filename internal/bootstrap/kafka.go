package bootstrap

import (
	"context"
	"fmt"

	"github.com/twmb/franz-go/pkg/kgo"
)

func NewKafkaClient(cfg KafkaConfig) (*kgo.Client, error) {
	client, err := kgo.NewClient(
		kgo.SeedBrokers(cfg.Brokers...),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize kafka client: %w", err)
	}

	// Ping to verify connection
	if err := client.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("kafka ping failed: %w", err)
	}

	return client, nil
}