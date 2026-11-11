package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/linkchat/api/internal/domain"
)

type AutomationRepository struct{ db *sql.DB }

func NewAutomationRepository(db *sql.DB) *AutomationRepository {
	return &AutomationRepository{db: db}
}

func (r *AutomationRepository) List(ctx context.Context, tenantID uint64) ([]domain.Automation, error) {
	rows, err := r.db.QueryContext(ctx,
		"SELECT id,tenant_id,name,COALESCE(description,''),event,conditions,actions,is_active,run_count,created_at,updated_at FROM automations WHERE tenant_id=? ORDER BY created_at DESC",
		tenantID)
	if err != nil {
		return nil, fmt.Errorf("automations.List: %w", err)
	}
	defer rows.Close()
	var list []domain.Automation
	for rows.Next() {
		var a domain.Automation
		var condJSON, actJSON sql.NullString
		if err = rows.Scan(&a.ID, &a.TenantID, &a.Name, &a.Description, &a.Event, &condJSON, &actJSON, &a.IsActive, &a.RunCount, &a.CreatedAt, &a.UpdatedAt); err != nil {
			return nil, err
		}
		if condJSON.Valid {
			json.Unmarshal([]byte(condJSON.String), &a.Conditions)
		}
		if actJSON.Valid {
			json.Unmarshal([]byte(actJSON.String), &a.Actions)
		}
		list = append(list, a)
	}
	return list, nil
}

func (r *AutomationRepository) FindByID(ctx context.Context, id, tenantID uint64) (*domain.Automation, error) {
	var a domain.Automation
	var condJSON, actJSON sql.NullString
	err := r.db.QueryRowContext(ctx,
		"SELECT id,tenant_id,name,COALESCE(description,''),event,conditions,actions,is_active,run_count,created_at,updated_at FROM automations WHERE id=? AND tenant_id=?",
		id, tenantID).Scan(&a.ID, &a.TenantID, &a.Name, &a.Description, &a.Event, &condJSON, &actJSON, &a.IsActive, &a.RunCount, &a.CreatedAt, &a.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if condJSON.Valid {
		json.Unmarshal([]byte(condJSON.String), &a.Conditions)
	}
	if actJSON.Valid {
		json.Unmarshal([]byte(actJSON.String), &a.Actions)
	}
	return &a, err
}

func (r *AutomationRepository) Create(ctx context.Context, a domain.Automation) (*domain.Automation, error) {
	condJSON, _ := json.Marshal(a.Conditions)
	actJSON, _ := json.Marshal(a.Actions)
	res, err := r.db.ExecContext(ctx,
		"INSERT INTO automations(tenant_id,name,description,event,conditions,actions) VALUES(?,?,?,?,?,?)",
		a.TenantID, a.Name, a.Description, a.Event, string(condJSON), string(actJSON))
	if err != nil {
		return nil, fmt.Errorf("automations.Create: %w", err)
	}
	id, _ := res.LastInsertId()
	return r.FindByID(ctx, uint64(id), a.TenantID)
}

func (r *AutomationRepository) ToggleActive(ctx context.Context, id, tenantID uint64, active bool) error {
	_, err := r.db.ExecContext(ctx, "UPDATE automations SET is_active=?,updated_at=NOW() WHERE id=? AND tenant_id=?", active, id, tenantID)
	return err
}

func (r *AutomationRepository) Delete(ctx context.Context, id, tenantID uint64) error {
	res, err := r.db.ExecContext(ctx, "DELETE FROM automations WHERE id=? AND tenant_id=?", id, tenantID)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return ErrNotFound
	}
	return nil
}
