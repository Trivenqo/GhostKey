package canonical

import (
	"time"
)

type IdentityType string

const (
	IdentityTypeIAMRole        IdentityType = "iam_role"
	IdentityTypeServiceAccount IdentityType = "service_account"
	IdentityTypeAPIKey         IdentityType = "api_key"
	IdentityTypeOAuthClient    IdentityType = "oauth_client"
	IdentityTypeSSHKey         IdentityType = "ssh_key"
	IdentityTypeTLSCertificate IdentityType = "tls_certificate"
	IdentityTypeCICDToken      IdentityType = "cicd_token"
	IdentityTypeAIAgent        IdentityType = "ai_agent"
	IdentityTypeMCPServer      IdentityType = "mcp_server"
)

type ProviderRef struct {
	Provider   string // e.g., "aws", "github", "kubernetes"
	ExternalID string // The provider's native ID (e.g., ARN)
	AccountRef string // AWS account ID, GitHub org, etc.
}

type Identity struct {
	ID            string
	ExternalRef   ProviderRef
	Type          IdentityType
	DisplayName   string
	CreatedAt     time.Time
	LastUsedAt    *time.Time
	Credentials   []Credential
	Permissions   []Permission
	Owner         *OwnerAssignment
	Relationships []Relationship
	Metadata      map[string]string // Opaque key-value for provider-specific extras
	RiskFactors   RiskFactorSet
}

type Credential struct {
	ID          string
	IdentityID  string
	Kind        string // "access_key", "cert", "token"
	CreatedAt   time.Time
	ExpiresAt   *time.Time
	LastRotated *time.Time
	Fingerprint string // Hash, NEVER the raw secret
}

type Permission struct {
	Action   string
	Resource string
	Effect   string // "allow" or "deny"
	Scope    string // "admin", "read", "write"
}

type OwnerAssignment struct {
	OwnerType  string // "user", "team", "service"
	OwnerRef   string
	Confidence ConfidenceScore
	Evidence   []Evidence
}

type Evidence struct {
	Source    string // "git_commit", "codeowners", "slack"
	Detail    string
	Timestamp time.Time
	Weight    float64
}

type Relationship struct {
	FromID string
	ToID   string
	Kind   string // "assumes", "creates", "calls"
}

type RiskFactorSet struct {
	Privilege         float64
	CredentialAgeDays int
	Unused            bool
	PublicExposure    bool
	InternetReachable bool
	AdminAccess       bool
	SecretsAccess     bool
	ProductionAccess  bool
	OwnerConfidence   float64
	BehaviorAnomaly   float64
	BlastRadius       float64
}

type ConfidenceScore float64 // 0.0 to 1.0

// NewConfidenceScore ensures the score stays within bounds.
func NewConfidenceScore(score float64) ConfidenceScore {
	if score < 0.0 {
		return 0.0
	}
	if score > 1.0 {
		return 1.0
	}
	return ConfidenceScore(score)
}