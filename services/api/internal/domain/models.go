package domain

import (
	"database/sql"
	"time"
)

// Tenant represents a workspace / company using LinkChat.
type Tenant struct {
	ID        uint64    `json:"id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TenantUser represents a user that belongs to a tenant.
type TenantUser struct {
	ID               uint64       `json:"id"`
	TenantID         uint64       `json:"tenant_id"`
	FullName         string       `json:"full_name"`
	Email            string       `json:"email"`
	Username         string       `json:"username"`
	PasswordHash     string       `json:"-"` // never serialized
	Role             string       `json:"role"`
	IsActive         bool         `json:"is_active"`
	EmailVerifiedAt  sql.NullTime `json:"-"`
	LastLoginAt      sql.NullTime `json:"-"`
	CreatedAt        time.Time    `json:"created_at"`
	UpdatedAt        time.Time    `json:"updated_at"`
}

// ─── Request / Response DTOs ────────────────────────────────

// RegisterRequest is the payload sent by the client to POST /api/v1/auth/register.
type RegisterRequest struct {
	WorkspaceName   string `json:"workspace_name"`
	FullName        string `json:"full_name"`
	Email           string `json:"email"`
	Username        string `json:"username"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

// LoginRequest is the payload sent by the client to POST /api/v1/auth/login.
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// MeResponse is what GET /api/v1/auth/me returns.
type MeResponse struct {
	UserID   uint64 `json:"user_id"`
	TenantID uint64 `json:"tenant_id"`
	FullName string `json:"full_name"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	Tenant   struct {
		ID   uint64 `json:"id"`
		Name string `json:"name"`
		Slug string `json:"slug"`
	} `json:"tenant"`
}

// SessionData is stored inside the secure server-side session cookie.
type SessionData struct {
	UserID   uint64
	TenantID uint64
	Role     string
}
