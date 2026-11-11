package domain

import "time"

type Automation struct {
	ID          uint64    `json:"id"`
	TenantID    uint64    `json:"tenant_id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	Event       string    `json:"event"`
	Conditions  any       `json:"conditions"`
	Actions     any       `json:"actions"`
	IsActive    bool      `json:"is_active"`
	RunCount    int       `json:"run_count"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateAutomationRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Event       string `json:"event"`
	Conditions  any    `json:"conditions"`
	Actions     any    `json:"actions"`
}
