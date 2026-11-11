package domain

import "time"

type Contact struct {
	ID         uint64    `json:"id"`
	TenantID   uint64    `json:"tenant_id"`
	Name       string    `json:"name"`
	Phone      string    `json:"phone"`
	Email      string    `json:"email,omitempty"`
	AvatarURL  string    `json:"avatar_url,omitempty"`
	Notes      string    `json:"notes,omitempty"`
	IsBlocked  bool      `json:"is_blocked"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type CreateContactRequest struct {
	Name  string `json:"name"`
	Phone string `json:"phone"`
	Email string `json:"email"`
	Notes string `json:"notes"`
}

type UpdateContactRequest struct {
	Name      string `json:"name"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
	Notes     string `json:"notes"`
	IsBlocked bool   `json:"is_blocked"`
}
