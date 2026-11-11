package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/linkchat/api/internal/domain"
)

// ErrNotFound is returned when a queried record does not exist.
var ErrNotFound = errors.New("record not found")

// ErrDuplicateUsername signals that the username is already taken.
var ErrDuplicateUsername = errors.New("username already exists")

// ErrDuplicateEmail signals that the email is already taken.
var ErrDuplicateEmail = errors.New("email already exists")

// AuthRepository handles all database operations related to authentication.
type AuthRepository struct {
	db *sql.DB
}

// NewAuthRepository creates a new AuthRepository.
func NewAuthRepository(db *sql.DB) *AuthRepository {
	return &AuthRepository{db: db}
}

// ─── Read ────────────────────────────────────────────────────

// FindUserByUsername looks up a tenant_user by username (case-insensitive via stored lowercase).
func (r *AuthRepository) FindUserByUsername(ctx context.Context, username string) (*domain.TenantUser, error) {
	const q = `
		SELECT id, tenant_id, full_name, email, username, password_hash,
		       role, is_active, email_verified_at, last_login_at, created_at, updated_at
		FROM tenant_users
		WHERE username = ?
		LIMIT 1`

	row := r.db.QueryRowContext(ctx, q, username)
	u, err := scanUser(row)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	return u, err
}

// FindUserByID looks up a tenant_user by its primary key.
func (r *AuthRepository) FindUserByID(ctx context.Context, userID, tenantID uint64) (*domain.TenantUser, error) {
	const q = `
		SELECT id, tenant_id, full_name, email, username, password_hash,
		       role, is_active, email_verified_at, last_login_at, created_at, updated_at
		FROM tenant_users
		WHERE id = ? AND tenant_id = ?
		LIMIT 1`

	row := r.db.QueryRowContext(ctx, q, userID, tenantID)
	u, err := scanUser(row)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	return u, err
}

// FindTenantByID retrieves a tenant by its PK.
func (r *AuthRepository) FindTenantByID(ctx context.Context, tenantID uint64) (*domain.Tenant, error) {
	const q = `
		SELECT id, name, slug, status, created_at, updated_at
		FROM tenants
		WHERE id = ?
		LIMIT 1`

	row := r.db.QueryRowContext(ctx, q, tenantID)
	var t domain.Tenant
	err := row.Scan(&t.ID, &t.Name, &t.Slug, &t.Status, &t.CreatedAt, &t.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	return &t, err
}

// UsernameExists checks if a username is already taken.
func (r *AuthRepository) UsernameExists(ctx context.Context, username string) (bool, error) {
	const q = `SELECT COUNT(1) FROM tenant_users WHERE username = ?`
	var count int
	err := r.db.QueryRowContext(ctx, q, username).Scan(&count)
	return count > 0, err
}

// EmailExists checks if an email is already taken.
func (r *AuthRepository) EmailExists(ctx context.Context, email string) (bool, error) {
	const q = `SELECT COUNT(1) FROM tenant_users WHERE email = ?`
	var count int
	err := r.db.QueryRowContext(ctx, q, email).Scan(&count)
	return count > 0, err
}

// ─── Write (transactional) ───────────────────────────────────

// CreateTenantAndAdmin creates a tenant and its first admin atomically.
// If either INSERT fails, the transaction is rolled back.
func (r *AuthRepository) CreateTenantAndAdmin(
	ctx context.Context,
	tenant domain.Tenant,
	user domain.TenantUser,
) (*domain.Tenant, *domain.TenantUser, error) {

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("CreateTenantAndAdmin: begin tx: %w", err)
	}

	// Rollback helper — only takes effect if tx was not committed.
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	// 1. Insert tenant
	const insertTenant = `
		INSERT INTO tenants (name, slug, status, created_at, updated_at)
		VALUES (?, ?, 'active', NOW(), NOW())`

	res, err := tx.ExecContext(ctx, insertTenant, tenant.Name, tenant.Slug)
	if err != nil {
		return nil, nil, fmt.Errorf("CreateTenantAndAdmin: insert tenant: %w", err)
	}
	tenantID, err := res.LastInsertId()
	if err != nil {
		return nil, nil, fmt.Errorf("CreateTenantAndAdmin: last insert tenant id: %w", err)
	}

	// 2. Insert admin user
	const insertUser = `
		INSERT INTO tenant_users
		    (tenant_id, full_name, email, username, password_hash, role, is_active, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, 'admin', 1, NOW(), NOW())`

	ures, err := tx.ExecContext(ctx, insertUser,
		tenantID, user.FullName, user.Email, user.Username, user.PasswordHash)
	if err != nil {
		return nil, nil, fmt.Errorf("CreateTenantAndAdmin: insert user: %w", err)
	}
	userID, err := ures.LastInsertId()
	if err != nil {
		return nil, nil, fmt.Errorf("CreateTenantAndAdmin: last insert user id: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return nil, nil, fmt.Errorf("CreateTenantAndAdmin: commit: %w", err)
	}

	tenant.ID = uint64(tenantID)
	user.ID = uint64(userID)
	user.TenantID = uint64(tenantID)

	return &tenant, &user, nil
}

// UpdateLastLogin sets last_login_at to NOW() for the given user.
func (r *AuthRepository) UpdateLastLogin(ctx context.Context, userID uint64) error {
	const q = `UPDATE tenant_users SET last_login_at = NOW() WHERE id = ?`
	_, err := r.db.ExecContext(ctx, q, userID)
	return err
}

// ─── Helpers ─────────────────────────────────────────────────

func scanUser(row *sql.Row) (*domain.TenantUser, error) {
	var u domain.TenantUser
	err := row.Scan(
		&u.ID, &u.TenantID, &u.FullName, &u.Email, &u.Username,
		&u.PasswordHash, &u.Role, &u.IsActive,
		&u.EmailVerifiedAt, &u.LastLoginAt,
		&u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &u, nil
}
