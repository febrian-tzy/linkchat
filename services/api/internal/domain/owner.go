package domain

import "time"

// Owner is the platform administrator — completely separate from tenant users.
type Owner struct {
	ID           uint64    `json:"id"`
	FullName     string    `json:"full_name"`
	Email        string    `json:"email"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"-"`
	IsActive     bool      `json:"is_active"`
	LastLoginAt  *time.Time `json:"last_login_at,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// OwnerLoginRequest is the payload for POST /api/v1/owner/auth/login.
type OwnerLoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// OwnerSessionData is stored inside the owner JWT.
type OwnerSessionData struct {
	OwnerID uint64
}

// OwnerMeResponse is returned by GET /api/v1/owner/auth/me.
type OwnerMeResponse struct {
	ID       uint64 `json:"id"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Role     string `json:"role"` // always "owner"
}
