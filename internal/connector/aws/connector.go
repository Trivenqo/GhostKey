package aws

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Trivenqo/GhostKey/internal/connector/sdk"
	"github.com/Trivenqo/GhostKey/internal/shared/canonical"
)

// Connector implements sdk.Connector for AWS IAM.
type Connector struct {
	isAuthenticated bool
	accountID       string
}

func NewConnector() *Connector {
	return &Connector{}
}

func (c *Connector) Name() string {
	return "aws"
}

// Authenticate validates the provided credentials (e.g., Access Key / Secret).
func (c *Connector) Authenticate(ctx context.Context, creds sdk.Credentials) error {
	accessKey, ok := creds["access_key_id"]
	if !ok || accessKey == "" {
		return errors.New("missing AWS access_key_id in credentials")
	}

	// In a real implementation, you would initialize the AWS SDK client here.
	// For this simulation, we'll pretend authentication succeeded.
	c.isAuthenticated = true
	c.accountID = "123456789012" // Mock AWS Account ID
	return nil
}

// Discover fetches IAM Users and Roles. We use a mock response to demonstrate pagination.
func (c *Connector) Discover(ctx context.Context, cursor sdk.Cursor) (sdk.Page, error) {
	if !c.isAuthenticated {
		return sdk.Page{}, errors.New("not authenticated")
	}

	// Simulated Pagination Logic
	var items []sdk.RawRecord
	var nextToken string
	var hasMore bool

	switch cursor.Token {
	case "":
		// First page: return an IAM User
		items = append(items, sdk.RawRecord{
			"resource_type": "iam_user",
			"arn":           fmt.Sprintf("arn:aws:iam::%s:user/alice.dev", c.accountID),
			"user_name":     "alice.dev",
			"create_date":   time.Now().Add(-8760 * time.Hour).Format(time.RFC3339), // 1 year ago
			"tags":          map[string]interface{}{"Team": "Engineering"},
		})
		nextToken = "page_2"
		hasMore = true

	case "page_2":
		// Second page: return an IAM Role
		items = append(items, sdk.RawRecord{
			"resource_type": "iam_role",
			"arn":           fmt.Sprintf("arn:aws:iam::%s:role/EKS-Worker-Node", c.accountID),
			"role_name":     "EKS-Worker-Node",
			"create_date":   time.Now().Add(-720 * time.Hour).Format(time.RFC3339), // 30 days ago
			"tags":          map[string]interface{}{"Environment": "Production"},
		})
		nextToken = ""
		hasMore = false

	default:
		return sdk.Page{}, fmt.Errorf("invalid pagination token: %s", cursor.Token)
	}

	return sdk.Page{
		Items: items,
		NextCursor: sdk.Cursor{
			Token:     nextToken,
			UpdatedAt: time.Now(),
		},
		HasMore: hasMore,
	}, nil
}

func (c *Connector) SupportsWebhook() bool {
	// AWS IAM doesn't support direct webhooks natively (usually requires EventBridge + SNS/SQS).
	// We rely on polling.
	return false
}

func (c *Connector) HandleWebhook(ctx context.Context, payload []byte) ([]canonical.Identity, error) {
	return nil, errors.New("webhooks not supported for AWS connector")
}

func (c *Connector) HealthCheck(ctx context.Context) sdk.HealthStatus {
	return sdk.HealthStatus{
		Healthy:    c.isAuthenticated,
		LastSyncAt: time.Now(),
		Message:    "AWS connector is operational",
	}
}