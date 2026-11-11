package service

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/crypto/bcrypt"

	"github.com/linkchat/api/internal/domain"
	"github.com/linkchat/api/internal/repository"
)

// Sentinel errors exposed to handler layer.
var (
	ErrUsernameTaken   = errors.New("username already taken")
	ErrEmailTaken      = errors.New("email already taken")
	ErrInvalidInput    = errors.New("invalid input")
	ErrInvalidCreds    = errors.New("invalid credentials")
	ErrUserInactive    = errors.New("user account is inactive")
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

// AuthService contains the business logic for authentication.
type AuthService struct {
	repo *repository.AuthRepository
}

// NewAuthService creates a new AuthService.
func NewAuthService(repo *repository.AuthRepository) *AuthService {
	return &AuthService{repo: repo}
}

// ─── Register ────────────────────────────────────────────────

// Register validates the request, checks uniqueness, hashes the password,
// and creates the tenant + admin atomically.
func (s *AuthService) Register(ctx context.Context, req domain.RegisterRequest) (*domain.Tenant, *domain.TenantUser, error) {
	// 1. Validate inputs
	if err := validateRegister(req); err != nil {
		return nil, nil, err
	}

	// 2. Normalize
	username := strings.ToLower(strings.TrimSpace(req.Username))
	email := strings.ToLower(strings.TrimSpace(req.Email))
	workspaceName := strings.TrimSpace(req.WorkspaceName)
	fullName := strings.TrimSpace(req.FullName)

	// 3. Uniqueness checks
	taken, err := s.repo.UsernameExists(ctx, username)
	if err != nil {
		return nil, nil, fmt.Errorf("register: check username: %w", err)
	}
	if taken {
		return nil, nil, ErrUsernameTaken
	}

	taken, err = s.repo.EmailExists(ctx, email)
	if err != nil {
		return nil, nil, fmt.Errorf("register: check email: %w", err)
	}
	if taken {
		return nil, nil, ErrEmailTaken
	}

	// 4. Hash password with bcrypt (cost 12 — secure & fast enough for web)
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
	if err != nil {
		return nil, nil, fmt.Errorf("register: hash password: %w", err)
	}

	// 5. Build domain objects
	tenant := domain.Tenant{
		Name:   workspaceName,
		Slug:   slugify(workspaceName),
		Status: "active",
	}
	user := domain.TenantUser{
		FullName:     fullName,
		Email:        email,
		Username:     username,
		PasswordHash: string(hash),
		Role:         "admin",
		IsActive:     true,
	}

	// 6. Atomic create (DB transaction inside repository)
	return s.repo.CreateTenantAndAdmin(ctx, tenant, user)
}

// ─── Login ───────────────────────────────────────────────────

// Login verifies credentials and returns the authenticated user.
func (s *AuthService) Login(ctx context.Context, req domain.LoginRequest) (*domain.TenantUser, error) {
	if strings.TrimSpace(req.Username) == "" || strings.TrimSpace(req.Password) == "" {
		return nil, ErrInvalidInput
	}

	username := strings.ToLower(strings.TrimSpace(req.Username))

	// 1. Find user
	user, err := s.repo.FindUserByUsername(ctx, username)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, ErrInvalidCreds // don't reveal whether username exists
	}
	if err != nil {
		return nil, fmt.Errorf("login: find user: %w", err)
	}

	// 2. Check account status
	if !user.IsActive {
		return nil, ErrUserInactive
	}

	// 3. Verify password
	if err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, ErrInvalidCreds
	}

	// 4. Update last_login_at (non-critical, ignore error)
	_ = s.repo.UpdateLastLogin(ctx, user.ID)

	return user, nil
}

// ─── Me ──────────────────────────────────────────────────────

// GetMe returns user and tenant data from a validated session.
// tenantID comes from the server-side session, never from client input.
func (s *AuthService) GetMe(ctx context.Context, userID, tenantID uint64) (*domain.MeResponse, error) {
	user, err := s.repo.FindUserByID(ctx, userID, tenantID)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, ErrInvalidCreds
	}
	if err != nil {
		return nil, fmt.Errorf("getMe: find user: %w", err)
	}

	tenant, err := s.repo.FindTenantByID(ctx, tenantID)
	if err != nil {
		return nil, fmt.Errorf("getMe: find tenant: %w", err)
	}

	me := &domain.MeResponse{
		UserID:   user.ID,
		TenantID: user.TenantID,
		FullName: user.FullName,
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role,
	}
	me.Tenant.ID = tenant.ID
	me.Tenant.Name = tenant.Name
	me.Tenant.Slug = tenant.Slug

	return me, nil
}

// ─── Helpers ─────────────────────────────────────────────────

func validateRegister(req domain.RegisterRequest) error {
	errs := []string{}

	if strings.TrimSpace(req.WorkspaceName) == "" {
		errs = append(errs, "workspace_name is required")
	}
	if strings.TrimSpace(req.FullName) == "" {
		errs = append(errs, "full_name is required")
	}
	if !emailRegex.MatchString(strings.TrimSpace(req.Email)) {
		errs = append(errs, "email is invalid")
	}
	if len(strings.TrimSpace(req.Username)) < 3 {
		errs = append(errs, "username must be at least 3 characters")
	}
	if len(req.Password) < 8 {
		errs = append(errs, "password must be at least 8 characters")
	}
	if req.Password != req.ConfirmPassword {
		errs = append(errs, "password and confirm_password do not match")
	}

	if len(errs) > 0 {
		return fmt.Errorf("%w: %s", ErrInvalidInput, strings.Join(errs, "; "))
	}
	return nil
}

// slugify converts a string to a URL-safe slug (lowercase, hyphens only).
func slugify(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	var b strings.Builder
	for _, r := range s {
		switch {
		case unicode.IsLetter(r) || unicode.IsDigit(r):
			b.WriteRune(r)
		case unicode.IsSpace(r) || r == '-' || r == '_':
			b.WriteRune('-')
		}
	}
	// Collapse multiple hyphens
	slug := regexp.MustCompile(`-{2,}`).ReplaceAllString(b.String(), "-")
	return strings.Trim(slug, "-")
}
