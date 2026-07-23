package infrastructure

import (
	"context"
	"encoding/json"
	"log"
	"regexp"
	"strings"

	"github.com/segmentio/kafka-go"

	"github.com/Trivenqo/GhostKey/internal/ownership/application"
)

type DiscoveredIdentityEvent struct {
	ID        string            `json:"id"` // Changed to string
	ARN       string            `json:"arn"`
	AccountID string            `json:"account_id"`
	Provider  string            `json:"provider"`
	Type      string            `json:"type"`
	Tags      map[string]string `json:"tags"`
}

type KafkaConsumer struct {
	reader  *kafka.Reader
	useCase *application.OwnershipUseCase
}

func NewKafkaConsumer(brokers []string, topic, groupID string, useCase *application.OwnershipUseCase) *KafkaConsumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		Topic:    topic,
		GroupID:  groupID,
		MinBytes: 10e3,
		MaxBytes: 10e6,
	})

	return &KafkaConsumer{
		reader:  reader,
		useCase: useCase,
	}
}

func (c *KafkaConsumer) Start(ctx context.Context) {
	log.Println("[Kafka Consumer] Subscribed to identity discovery events...")
	for {
		select {
		case <-ctx.Done():
			c.reader.Close()
			return
		default:
			msg, err := c.reader.ReadMessage(ctx)
			if err != nil {
				if ctx.Err() != nil {
					return
				}
				log.Printf("Error reading message from Kafka: %v", err)
				continue
			}

			var event DiscoveredIdentityEvent
			if err := json.Unmarshal(msg.Value, &event); err != nil {
				log.Printf("Failed to unmarshal discovery event: %v", err)
				continue
			}

			c.processEvent(ctx, event)
		}
	}
}

func (c *KafkaConsumer) processEvent(ctx context.Context, event DiscoveredIdentityEvent) {
	ownerEmail, teamName, department := extractContext(event)

	if ownerEmail == "" && teamName == "" && department == "" {
		return
	}

	err := c.useCase.AssignAutoOwnership(ctx, event.ID, ownerEmail, teamName, department)
	if err != nil {
		log.Printf("Failed to auto-assign ownership for identity %s: %v", event.ID, err)
		return
	}

	log.Printf("[Ownership Consumer] Successfully mapped ownership for identity %s (Team: %s, Owner: %s)", event.ID, teamName, ownerEmail)
}

func extractContext(event DiscoveredIdentityEvent) (ownerEmail, teamName, department string) {
	for k, v := range event.Tags {
		lowerK := strings.ToLower(k)
		switch lowerK {
		case "owner", "owner_email", "email":
			ownerEmail = v
		case "team", "team_name":
			teamName = v
		case "department", "dept":
			department = v
		}
	}

	if ownerEmail == "" {
		re := regexp.MustCompile(`user/(?:([^/]+)/)?([^/]+)`)
		matches := re.FindStringSubmatch(event.ARN)
		if len(matches) > 2 {
			if teamName == "" && matches[1] != "" {
				teamName = matches[1]
			}
			username := matches[2]
			if strings.Contains(username, "@") {
				ownerEmail = username
			}
		}
	}

	return ownerEmail, teamName, department
}