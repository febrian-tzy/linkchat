package domain

import "time"

type Webhook struct {
	ID        uint64    `json:"id"`
	TenantID  uint64    `json:"tenant_id"`
	URL       string    `json:"url"`
	Secret    string    `json:"secret,omitempty"`
	Events    []string  `json:"events"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateWebhookRequest struct {
	URL    string   `json:"url"`
	Secret string   `json:"secret"`
	Events []string `json:"events"`
}

type APIKey struct {
	ID         uint64    `json:"id"`
	TenantID   uint64    `json:"tenant_id"`
	Name       string    `json:"name"`
	KeyHash    string    `json:"-"`
	KeyPrefix  string    `json:"key_prefix"`
	IsActive   bool      `json:"is_active"`
	CreatedAt  time.Time `json:"created_at"`
}

type CreateAPIKeyRequest struct {
	Name string `json:"name"`
}

type CreateAPIKeyResponse struct {
	APIKey
	PlainKey string `json:"key"` // only returned once
}
