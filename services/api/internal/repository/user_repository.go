package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/linkchat/api/internal/domain"
)

type UserRepository struct{ db *sql.DB }

func NewUserRepository(db *sql.DB) *UserRepository { return &UserRepository{db: db} }

func (r *UserRepository) List(ctx context.Context, tenantID uint64) ([]domain.TenantUser, error) {
	rows, err := r.db.QueryContext(ctx,
		"SELECT id,tenant_id,full_name,email,username,role,is_active,created_at,updated_at FROM tenant_users WHERE tenant_id=? ORDER BY created_at ASC",
		tenantID)
	if err != nil {
		return nil, fmt.Errorf("users.List: %w", err)
	}
	defer rows.Close()
	var users []domain.TenantUser
	for rows.Next() {
		var u domain.TenantUser
		if err = rows.Scan(&u.ID, &u.TenantID, &u.FullName, &u.Email, &u.Username, &u.Role, &u.IsActive, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func (r *UserRepository) FindByID(ctx context.Context, id, tenantID uint64) (*domain.TenantUser, error) {
	var u domain.TenantUser
	err := r.db.QueryRowContext(ctx,
		"SELECT id,tenant_id,full_name,email,username,password_hash,role,is_active,created_at,updated_at FROM tenant_users WHERE id=? AND tenant_id=?",
		id, tenantID).Scan(&u.ID, &u.TenantID, &u.FullName, &u.Email, &u.Username, &u.PasswordHash, &u.Role, &u.IsActive, &u.CreatedAt, &u.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	return &u, err
}

func (r *UserRepository) Create(ctx context.Context, u domain.TenantUser) (*domain.TenantUser, error) {
	res, err := r.db.ExecContext(ctx,
		"INSERT INTO tenant_users(tenant_id,full_name,email,username,password_hash,role,is_active) VALUES(?,?,?,?,?,?,1)",
		u.TenantID, u.FullName, u.Email, u.Username, u.PasswordHash, u.Role)
	if err != nil {
		return nil, fmt.Errorf("users.Create: %w", err)
	}
	id, _ := res.LastInsertId()
	return r.FindByID(ctx, uint64(id), u.TenantID)
}

func (r *UserRepository) UpdateProfile(ctx context.Context, id, tenantID uint64, fullName, email string) error {
	_, err := r.db.ExecContext(ctx,
		"UPDATE tenant_users SET full_name=?,email=?,updated_at=NOW() WHERE id=? AND tenant_id=?",
		fullName, email, id, tenantID)
	return err
}

func (r *UserRepository) UpdatePassword(ctx context.Context, id uint64, hash string) error {
	_, err := r.db.ExecContext(ctx, "UPDATE tenant_users SET password_hash=?,updated_at=NOW() WHERE id=?", hash, id)
	return err
}

func (r *UserRepository) ToggleActive(ctx context.Context, id, tenantID uint64, active bool) error {
	_, err := r.db.ExecContext(ctx, "UPDATE tenant_users SET is_active=?,updated_at=NOW() WHERE id=? AND tenant_id=?", active, id, tenantID)
	return err
}

func (r *UserRepository) Delete(ctx context.Context, id, tenantID uint64) error {
	res, err := r.db.ExecContext(ctx, "DELETE FROM tenant_users WHERE id=? AND tenant_id=?", id, tenantID)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return ErrNotFound
	}
	return nil
}
