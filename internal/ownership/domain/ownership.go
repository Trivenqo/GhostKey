package domain

import (
	"time"

	"github.com/google/uuid"
)

type MappingSource string

const (
	SourceAuto   MappingSource = "AUTO"
	SourceManual MappingSource = "MANUAL"
)

type Ownership struct {
	ID            uuid.UUID     `json:"id"`
	IdentityID    string        `json:"identity_id"` // Changed to string
	OwnerEmail    string        `json:"owner_email,omitempty"`
	TeamName      string        `json:"team_name,omitempty"`
	Department    string        `json:"department,omitempty"`
	MappingSource MappingSource `json:"mapping_source"`
	CreatedAt     time.Time     `json:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at"`
}

type IdentityOwnershipDTO struct {
	IdentityID    string        `json:"identity_id"` // Changed to string
	ARN           string        `json:"arn"`
	AccountID     string        `json:"account_id"`
	Provider      string        `json:"provider"`
	Type          string        `json:"type"`
	OwnerEmail    string        `json:"owner_email,omitempty"`
	TeamName      string        `json:"team_name,omitempty"`
	Department    string        `json:"department,omitempty"`
	MappingSource MappingSource `json:"mapping_source,omitempty"`
	IsMapped      bool          `json:"is_mapped"`
}