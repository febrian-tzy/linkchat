package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/linkchat/api/internal/domain"
)

type ContactRepository struct{ db *sql.DB }

func NewContactRepository(db *sql.DB) *ContactRepository {
	return &ContactRepository{db: db}
}

func (r *ContactRepository) List(ctx context.Context, tenantID uint64, search string, limit, offset int) ([]domain.Contact, int, error) {
	where := "tenant_id = ?"
	args := []any{tenantID}
	if search != "" {
		where += " AND (name LIKE ? OR phone LIKE ? OR email LIKE ?)"
		s := "%" + search + "%"
		args = append(args, s, s, s)
	}

	var total int
	err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM contacts WHERE "+where, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("contacts.List count: %w", err)
	}

	args = append(args, limit, offset)
	rows, err := r.db.QueryContext(ctx,
		"SELECT id,tenant_id,name,phone,COALESCE(email,''),COALESCE(avatar_url,''),COALESCE(notes,''),is_blocked,created_at,updated_at FROM contacts WHERE "+where+" ORDER BY created_at DESC LIMIT ? OFFSET ?",
		args...)
	if err != nil {
		return nil, 0, fmt.Errorf("contacts.List: %w", err)
	}
	defer rows.Close()

	var contacts []domain.Contact
	for rows.Next() {
		var c domain.Contact
		if err = rows.Scan(&c.ID, &c.TenantID, &c.Name, &c.Phone, &c.Email, &c.AvatarURL, &c.Notes, &c.IsBlocked, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, 0, err
		}
		contacts = append(contacts, c)
	}
	return contacts, total, nil
}

func (r *ContactRepository) FindByID(ctx context.Context, id, tenantID uint64) (*domain.Contact, error) {
	row := r.db.QueryRowContext(ctx,
		"SELECT id,tenant_id,name,phone,COALESCE(email,''),COALESCE(avatar_url,''),COALESCE(notes,''),is_blocked,created_at,updated_at FROM contacts WHERE id=? AND tenant_id=?",
		id, tenantID)
	var c domain.Contact
	err := row.Scan(&c.ID, &c.TenantID, &c.Name, &c.Phone, &c.Email, &c.AvatarURL, &c.Notes, &c.IsBlocked, &c.CreatedAt, &c.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	return &c, err
}

func (r *ContactRepository) Create(ctx context.Context, c domain.Contact) (*domain.Contact, error) {
	res, err := r.db.ExecContext(ctx,
		"INSERT INTO contacts(tenant_id,name,phone,email,notes) VALUES(?,?,?,?,?)",
		c.TenantID, c.Name, c.Phone, c.Email, c.Notes)
	if err != nil {
		return nil, fmt.Errorf("contacts.Create: %w", err)
	}
	id, _ := res.LastInsertId()
	return r.FindByID(ctx, uint64(id), c.TenantID)
}

func (r *ContactRepository) Update(ctx context.Context, c domain.Contact) (*domain.Contact, error) {
	_, err := r.db.ExecContext(ctx,
		"UPDATE contacts SET name=?,phone=?,email=?,notes=?,is_blocked=?,updated_at=NOW() WHERE id=? AND tenant_id=?",
		c.Name, c.Phone, c.Email, c.Notes, c.IsBlocked, c.ID, c.TenantID)
	if err != nil {
		return nil, fmt.Errorf("contacts.Update: %w", err)
	}
	return r.FindByID(ctx, c.ID, c.TenantID)
}

func (r *ContactRepository) Delete(ctx context.Context, id, tenantID uint64) error {
	res, err := r.db.ExecContext(ctx, "DELETE FROM contacts WHERE id=? AND tenant_id=?", id, tenantID)
	if err != nil {
		return fmt.Errorf("contacts.Delete: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return ErrNotFound
	}
	return nil
}
