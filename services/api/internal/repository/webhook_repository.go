package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/linkchat/api/internal/domain"
)

type WebhookRepository struct{ db *sql.DB }

func NewWebhookRepository(db *sql.DB) *WebhookRepository { return &WebhookRepository{db: db} }

func (r *WebhookRepository) List(ctx context.Context, tenantID uint64) ([]domain.Webhook, error) {
	rows, err := r.db.QueryContext(ctx,
		"SELECT id,tenant_id,url,COALESCE(secret,''),events,is_active,created_at,updated_at FROM webhooks WHERE tenant_id=? ORDER BY created_at DESC",
		tenantID)
	if err != nil {
		return nil, fmt.Errorf("webhooks.List: %w", err)
	}
	defer rows.Close()
	var list []domain.Webhook
	for rows.Next() {
		var w domain.Webhook
		var eventsJSON string
		if err = rows.Scan(&w.ID, &w.TenantID, &w.URL, &w.Secret, &eventsJSON, &w.IsActive, &w.CreatedAt, &w.UpdatedAt); err != nil {
			return nil, err
		}
		json.Unmarshal([]byte(eventsJSON), &w.Events)
		list = append(list, w)
	}
	return list, nil
}

func (r *WebhookRepository) FindByID(ctx context.Context, id, tenantID uint64) (*domain.Webhook, error) {
	var w domain.Webhook
	var eventsJSON string
	err := r.db.QueryRowContext(ctx,
		"SELECT id,tenant_id,url,COALESCE(secret,''),events,is_active,created_at,updated_at FROM webhooks WHERE id=? AND tenant_id=?",
		id, tenantID).Scan(&w.ID, &w.TenantID, &w.URL, &w.Secret, &eventsJSON, &w.IsActive, &w.CreatedAt, &w.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	json.Unmarshal([]byte(eventsJSON), &w.Events)
	return &w, err
}

func (r *WebhookRepository) Create(ctx context.Context, w domain.Webhook) (*domain.Webhook, error) {
	eventsJSON, _ := json.Marshal(w.Events)
	res, err := r.db.ExecContext(ctx,
		"INSERT INTO webhooks(tenant_id,url,secret,events) VALUES(?,?,?,?)",
		w.TenantID, w.URL, w.Secret, string(eventsJSON))
	if err != nil {
		return nil, fmt.Errorf("webhooks.Create: %w", err)
	}
	id, _ := res.LastInsertId()
	return r.FindByID(ctx, uint64(id), w.TenantID)
}

func (r *WebhookRepository) Delete(ctx context.Context, id, tenantID uint64) error {
	res, err := r.db.ExecContext(ctx, "DELETE FROM webhooks WHERE id=? AND tenant_id=?", id, tenantID)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return ErrNotFound
	}
	return nil
}

// ── API Keys ─────────────────────────────────────────────────

type APIKeyRepository struct{ db *sql.DB }

func NewAPIKeyRepository(db *sql.DB) *APIKeyRepository { return &APIKeyRepository{db: db} }

func (r *APIKeyRepository) List(ctx context.Context, tenantID uint64) ([]domain.APIKey, error) {
	rows, err := r.db.QueryContext(ctx,
		"SELECT id,tenant_id,name,key_prefix,is_active,created_at FROM api_keys WHERE tenant_id=? ORDER BY created_at DESC",
		tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []domain.APIKey
	for rows.Next() {
		var k domain.APIKey
		rows.Scan(&k.ID, &k.TenantID, &k.Name, &k.KeyPrefix, &k.IsActive, &k.CreatedAt)
		list = append(list, k)
	}
	return list, nil
}

func (r *APIKeyRepository) Create(ctx context.Context, k domain.APIKey) (*domain.APIKey, error) {
	res, err := r.db.ExecContext(ctx,
		"INSERT INTO api_keys(tenant_id,name,key_hash,key_prefix) VALUES(?,?,?,?)",
		k.TenantID, k.Name, k.KeyHash, k.KeyPrefix)
	if err != nil {
		return nil, err
	}
	id, _ := res.LastInsertId()
	k.ID = uint64(id)
	return &k, nil
}

func (r *APIKeyRepository) Delete(ctx context.Context, id, tenantID uint64) error {
	res, err := r.db.ExecContext(ctx, "DELETE FROM api_keys WHERE id=? AND tenant_id=?", id, tenantID)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return ErrNotFound
	}
	return nil
}
