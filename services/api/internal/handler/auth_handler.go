package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/linkchat/api/internal/domain"
	"github.com/linkchat/api/internal/middleware"
	"github.com/linkchat/api/internal/service"
	"github.com/linkchat/api/internal/session"
)

// AuthHandler handles all /api/v1/auth/* routes.
type AuthHandler struct {
	svc   *service.AuthService
	store *session.Store
}

// NewAuthHandler creates a new AuthHandler.
func NewAuthHandler(svc *service.AuthService, store *session.Store) *AuthHandler {
	return &AuthHandler{svc: svc, store: store}
}

// ─── POST /api/v1/auth/register ──────────────────────────────

// Register creates a new tenant workspace and its first admin user.
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req domain.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	tenant, user, err := h.svc.Register(r.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrUsernameTaken):
			jsonError(w, "username already taken", http.StatusConflict)
		case errors.Is(err, service.ErrEmailTaken):
			jsonError(w, "email already taken", http.StatusConflict)
		case errors.Is(err, service.ErrInvalidInput):
			jsonError(w, err.Error(), http.StatusUnprocessableEntity)
		default:
			jsonError(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}

	jsonResponse(w, http.StatusCreated, map[string]any{
		"message": "registration successful",
		"tenant": map[string]any{
			"id":   tenant.ID,
			"name": tenant.Name,
			"slug": tenant.Slug,
		},
		"user": map[string]any{
			"id":       user.ID,
			"username": user.Username,
			"role":     user.Role,
		},
	})
}

// ─── POST /api/v1/auth/login ─────────────────────────────────

// Login authenticates a tenant user and creates a server-side session.
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req domain.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	user, err := h.svc.Login(r.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidCreds):
			jsonError(w, "invalid username or password", http.StatusUnauthorized)
		case errors.Is(err, service.ErrUserInactive):
			jsonError(w, "account is inactive", http.StatusForbidden)
		case errors.Is(err, service.ErrInvalidInput):
			jsonError(w, "username and password are required", http.StatusUnprocessableEntity)
		default:
			jsonError(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}

	// Create JWT token; TenantID comes from the DB, never from client input.
	sd := domain.SessionData{
		UserID:   user.ID,
		TenantID: user.TenantID,
		Role:     user.Role,
	}
	token, err := h.store.CreateToken(sd)
	if err != nil {
		jsonError(w, "could not create session", http.StatusInternalServerError)
		return
	}

	jsonResponse(w, http.StatusOK, map[string]any{
		"message": "login successful",
		"token":   token,
	})
}

// ─── POST /api/v1/auth/logout ────────────────────────────────

// Logout destroys the session cookie.
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	if err := h.store.Clear(w, r); err != nil {
		jsonError(w, "could not clear session", http.StatusInternalServerError)
		return
	}
	jsonResponse(w, http.StatusOK, map[string]any{"message": "logged out"})
}

// ─── GET /api/v1/auth/me ─────────────────────────────────────

// Me returns the authenticated user and tenant info.
// TenantID is sourced from the server-side session — never from client input.
func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	sd := middleware.GetSession(r)
	if sd == nil {
		jsonError(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	me, err := h.svc.GetMe(r.Context(), sd.UserID, sd.TenantID)
	if err != nil {
		jsonError(w, "internal server error", http.StatusInternalServerError)
		return
	}

	jsonResponse(w, http.StatusOK, me)
}

// ─── Helpers ─────────────────────────────────────────────────

func jsonResponse(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func jsonError(w http.ResponseWriter, msg string, status int) {
	jsonResponse(w, status, map[string]string{"error": msg})
}
