package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/linkchat/api/internal/domain"
)

type BroadcastRepository struct{ db *sql.DB }

func NewBroadcastRepository(db *sql.DB) *BroadcastRepository {
	return &BroadcastRepository{db: db}
}

func (r *BroadcastRepository) List(ctx context.Context, tenantID uint64) ([]domain.Broadcast, error) {
	rows, err := r.db.QueryContext(ctx,
		"SELECT id,tenant_id,inbox_id,name,message,status,total,sent_count,failed_count,created_by,created_at,updated_at FROM broadcasts WHERE tenant_id=? ORDER BY created_at DESC",
		tenantID)
	if err != nil {
		return nil, fmt.Errorf("broadcasts.List: %w", err)
	}
	defer rows.Close()
	var list []domain.Broadcast
	for rows.Next() {
		var b domain.Broadcast
		if err = rows.Scan(&b.ID, &b.TenantID, &b.InboxID, &b.Name, &b.Message, &b.Status,
			&b.Total, &b.SentCount, &b.FailedCount, &b.CreatedBy, &b.CreatedAt, &b.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, b)
	}
	return list, nil
}

func (r *BroadcastRepository) FindByID(ctx context.Context, id, tenantID uint64) (*domain.Broadcast, error) {
	var b domain.Broadcast
	err := r.db.QueryRowContext(ctx,
		"SELECT id,tenant_id,inbox_id,name,message,status,total,sent_count,failed_count,created_by,created_at,updated_at FROM broadcasts WHERE id=? AND tenant_id=?",
		id, tenantID).Scan(&b.ID, &b.TenantID, &b.InboxID, &b.Name, &b.Message, &b.Status,
		&b.Total, &b.SentCount, &b.FailedCount, &b.CreatedBy, &b.CreatedAt, &b.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	return &b, err
}

func (r *BroadcastRepository) Create(ctx context.Context, b domain.Broadcast, contactIDs []uint64) (*domain.Broadcast, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	res, err := tx.ExecContext(ctx,
		"INSERT INTO broadcasts(tenant_id,inbox_id,name,message,status,total,created_by) VALUES(?,?,?,?,'draft',?,?)",
		b.TenantID, b.InboxID, b.Name, b.Message, len(contactIDs), b.CreatedBy)
	if err != nil {
		return nil, fmt.Errorf("broadcasts.Create: %w", err)
	}
	id, _ := res.LastInsertId()

	for _, cid := range contactIDs {
		if _, err = tx.ExecContext(ctx,
			"INSERT INTO broadcast_recipients(broadcast_id,contact_id) VALUES(?,?)", id, cid); err != nil {
			return nil, fmt.Errorf("broadcasts.Create recipients: %w", err)
		}
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return r.FindByID(ctx, uint64(id), b.TenantID)
}

func (r *BroadcastRepository) Delete(ctx context.Context, id, tenantID uint64) error {
	res, err := r.db.ExecContext(ctx, "DELETE FROM broadcasts WHERE id=? AND tenant_id=?", id, tenantID)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return ErrNotFound
	}
	return nil
}
