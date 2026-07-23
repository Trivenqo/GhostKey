package kafka

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/Trivenqo/GhostKey/internal/shared/canonical"
)

const (
	TopicIdentityDiscovered = "ghostkey.discovery.identity_discovered"
)

type Publisher struct {
	client *kgo.Client
}

func NewPublisher(client *kgo.Client) *Publisher {
	return &Publisher{client: client}
}

func (p *Publisher) PublishIdentityDiscovered(ctx context.Context, identity canonical.Identity) error {
	payload, err := json.Marshal(identity)
	if err != nil {
		return fmt.Errorf("failed to marshal identity event: %w", err)
	}

	record := &kgo.Record{
		Topic: TopicIdentityDiscovered,
		Value: payload,
		// Using the external ID as the routing key ensures events for 
		// the same identity are processed in order by consumers
		Key:   []byte(identity.ExternalRef.ExternalID),
	}

	// ProduceSync blocks until acknowledged by the broker to guarantee delivery
	if err := p.client.ProduceSync(ctx, record).FirstErr(); err != nil {
		return fmt.Errorf("failed to publish to kafka: %w", err)
	}

	return nil
} 