package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/linkchat/api/internal/domain"
)

// OwnerRepository handles DB operations for platform owners.
type OwnerRepository struct{ db *sql.DB }

func NewOwnerRepository(db *sql.DB) *OwnerRepository {
	return &OwnerRepository{db: db}
}

func (r *OwnerRepository) FindByUsername(ctx context.Context, username string) (*domain.Owner, error) {
	var o domain.Owner
	var lastLogin sql.NullTime
	err := r.db.QueryRowContext(ctx,
		"SELECT id,full_name,email,username,password_hash,is_active,last_login_at,created_at,updated_at FROM owners WHERE username=? LIMIT 1",
		username,
	).Scan(&o.ID, &o.FullName, &o.Email, &o.Username, &o.PasswordHash, &o.IsActive, &lastLogin, &o.CreatedAt, &o.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("owner.FindByUsername: %w", err)
	}
	if lastLogin.Valid {
		o.LastLoginAt = &lastLogin.Time
	}
	return &o, nil
}

func (r *OwnerRepository) FindByID(ctx context.Context, id uint64) (*domain.Owner, error) {
	var o domain.Owner
	var lastLogin sql.NullTime
	err := r.db.QueryRowContext(ctx,
		"SELECT id,full_name,email,username,password_hash,is_active,last_login_at,created_at,updated_at FROM owners WHERE id=? LIMIT 1",
		id,
	).Scan(&o.ID, &o.FullName, &o.Email, &o.Username, &o.PasswordHash, &o.IsActive, &lastLogin, &o.CreatedAt, &o.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("owner.FindByID: %w", err)
	}
	if lastLogin.Valid {
		o.LastLoginAt = &lastLogin.Time
	}
	return &o, nil
}

func (r *OwnerRepository) UpdateLastLogin(ctx context.Context, id uint64) error {
	_, err := r.db.ExecContext(ctx, "UPDATE owners SET last_login_at=NOW() WHERE id=?", id)
	return err
}

// ListTenants returns all tenants for the owner dashboard.
func (r *OwnerRepository) ListTenants(ctx context.Context) ([]domain.Tenant, error) {
	rows, err := r.db.QueryContext(ctx,
		"SELECT id,name,slug,status,created_at,updated_at FROM tenants ORDER BY created_at DESC")
	if err != nil {
		return nil, fmt.Errorf("owner.ListTenants: %w", err)
	}
	defer rows.Close()
	var tenants []domain.Tenant
	for rows.Next() {
		var t domain.Tenant
		if err = rows.Scan(&t.ID, &t.Name, &t.Slug, &t.Status, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		tenants = append(tenants, t)
	}
	return tenants, nil
}

// PlatformStats returns aggregate numbers for the owner dashboard.
func (r *OwnerRepository) PlatformStats(ctx context.Context) (map[string]int, error) {
	stats := map[string]int{}
	queries := map[string]string{
		"total_tenants":       "SELECT COUNT(*) FROM tenants",
		"active_tenants":      "SELECT COUNT(*) FROM tenants WHERE status='active'",
		"total_users":         "SELECT COUNT(*) FROM tenant_users",
		"total_contacts":      "SELECT COUNT(*) FROM contacts",
		"total_conversations": "SELECT COUNT(*) FROM conversations",
	}
	for key, q := range queries {
		var n int
		r.db.QueryRowContext(ctx, q).Scan(&n)
		stats[key] = n
	}
	return stats, nil
}
