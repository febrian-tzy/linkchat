package domain

import (
	"database/sql"
	"time"
)

type Conversation struct {
	ID             uint64       `json:"id"`
	TenantID       uint64       `json:"tenant_id"`
	InboxID        uint64       `json:"inbox_id"`
	ContactID      uint64       `json:"contact_id"`
	AssignedUserID sql.NullInt64 `json:"-"`
	AssignedUser   *string      `json:"assigned_user,omitempty"`
	Status         string       `json:"status"`
	UnreadCount    int          `json:"unread_count"`
	LastMessageAt  sql.NullTime `json:"-"`
	LastMessageAtV *time.Time   `json:"last_message_at,omitempty"`
	Contact        *Contact     `json:"contact,omitempty"`
	Inbox          *Inbox       `json:"inbox,omitempty"`
	CreatedAt      time.Time    `json:"created_at"`
	UpdatedAt      time.Time    `json:"updated_at"`
}

type Message struct {
	ID             uint64    `json:"id"`
	ConversationID uint64    `json:"conversation_id"`
	TenantID       uint64    `json:"tenant_id"`
	SenderType     string    `json:"sender_type"`
	SenderID       *uint64   `json:"sender_id,omitempty"`
	Content        string    `json:"content"`
	ContentType    string    `json:"content_type"`
	MediaURL       string    `json:"media_url,omitempty"`
	IsRead         bool      `json:"is_read"`
	CreatedAt      time.Time `json:"created_at"`
}

type CreateConversationRequest struct {
	InboxID   uint64 `json:"inbox_id"`
	ContactID uint64 `json:"contact_id"`
}

type SendMessageRequest struct {
	Content     string `json:"content"`
	ContentType string `json:"content_type"`
}

type UpdateConversationRequest struct {
	Status         string  `json:"status"`
	AssignedUserID *uint64 `json:"assigned_user_id"`
}
