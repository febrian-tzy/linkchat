package domain

import (
	"database/sql"
	"time"
)

type Broadcast struct {
	ID          uint64       `json:"id"`
	TenantID    uint64       `json:"tenant_id"`
	InboxID     uint64       `json:"inbox_id"`
	Name        string       `json:"name"`
	Message     string       `json:"message"`
	Status      string       `json:"status"`
	ScheduledAt sql.NullTime `json:"-"`
	SentAt      sql.NullTime `json:"-"`
	Total       int          `json:"total"`
	SentCount   int          `json:"sent_count"`
	FailedCount int          `json:"failed_count"`
	CreatedBy   uint64       `json:"created_by"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}

type CreateBroadcastRequest struct {
	InboxID      uint64   `json:"inbox_id"`
	Name         string   `json:"name"`
	Message      string   `json:"message"`
	ContactIDs   []uint64 `json:"contact_ids"`
	ScheduledAt  string   `json:"scheduled_at"`
}
