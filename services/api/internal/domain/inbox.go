package domain

import (
	"database/sql"
	"time"
)

type Inbox struct {
	ID          uint64         `json:"id"`
	TenantID    uint64         `json:"tenant_id"`
	Name        string         `json:"name"`
	ChannelType string         `json:"channel_type"`
	IsEnabled   bool           `json:"is_enabled"`
	Config      sql.NullString `json:"-"`
	ConfigJSON  any            `json:"config,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

type CreateInboxRequest struct {
	Name        string `json:"name"`
	ChannelType string `json:"channel_type"`
}
