package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/linkchat/api/internal/domain"
)

type InboxRepository struct{ db *sql.DB }

func NewInboxRepository(db *sql.DB) *InboxRepository { return &InboxRepository{db: db} }

func (r *InboxRepository) List(ctx context.Context, tenantID uint64) ([]domain.Inbox, error) {
	rows, err := r.db.QueryContext(ctx,
		"SELECT id,tenant_id,name,channel_type,is_enabled,created_at,updated_at FROM inboxes WHERE tenant_id=? ORDER BY created_at DESC",
		tenantID)
	if err != nil {
		return nil, fmt.Errorf("inboxes.List: %w", err)
	}
	defer rows.Close()
	var inboxes []domain.Inbox
	for rows.Next() {
		var i domain.Inbox
		if err = rows.Scan(&i.ID, &i.TenantID, &i.Name, &i.ChannelType, &i.IsEnabled, &i.CreatedAt, &i.UpdatedAt); err != nil {
			return nil, err
		}
		inboxes = append(inboxes, i)
	}
	return inboxes, nil
}

func (r *InboxRepository) FindByID(ctx context.Context, id, tenantID uint64) (*domain.Inbox, error) {
	var i domain.Inbox
	err := r.db.QueryRowContext(ctx,
		"SELECT id,tenant_id,name,channel_type,is_enabled,created_at,updated_at FROM inboxes WHERE id=? AND tenant_id=?",
		id, tenantID).Scan(&i.ID, &i.TenantID, &i.Name, &i.ChannelType, &i.IsEnabled, &i.CreatedAt, &i.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	return &i, err
}

func (r *InboxRepository) Create(ctx context.Context, i domain.Inbox) (*domain.Inbox, error) {
	res, err := r.db.ExecContext(ctx,
		"INSERT INTO inboxes(tenant_id,name,channel_type) VALUES(?,?,?)",
		i.TenantID, i.Name, i.ChannelType)
	if err != nil {
		return nil, fmt.Errorf("inboxes.Create: %w", err)
	}
	id, _ := res.LastInsertId()
	return r.FindByID(ctx, uint64(id), i.TenantID)
}

func (r *InboxRepository) Delete(ctx context.Context, id, tenantID uint64) error {
	res, err := r.db.ExecContext(ctx, "DELETE FROM inboxes WHERE id=? AND tenant_id=?", id, tenantID)
	if err != nil {
		return fmt.Errorf("inboxes.Delete: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return ErrNotFound
	}
	return nil
}
