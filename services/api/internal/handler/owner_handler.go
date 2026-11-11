package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"github.com/linkchat/api/internal/domain"
	"github.com/linkchat/api/internal/middleware"
	"github.com/linkchat/api/internal/repository"
	"github.com/linkchat/api/internal/session"
)

// OwnerHandler handles all /api/v1/owner/* routes.
// Completely separate from AuthHandler (tenant auth).
type OwnerHandler struct {
	repo  *repository.OwnerRepository
	store *session.OwnerStore
}

func NewOwnerHandler(repo *repository.OwnerRepository, store *session.OwnerStore) *OwnerHandler {
	return &OwnerHandler{repo: repo, store: store}
}

// POST /api/v1/owner/auth/login
func (h *OwnerHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req domain.OwnerLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	req.Username = strings.ToLower(strings.TrimSpace(req.Username))
	if req.Username == "" || req.Password == "" {
		jsonError(w, "username and password are required", http.StatusUnprocessableEntity)
		return
	}

	owner, err := h.repo.FindByUsername(r.Context(), req.Username)
	if errors.Is(err, repository.ErrNotFound) {
		jsonError(w, "invalid username or password", http.StatusUnauthorized)
		return
	}
	if err != nil {
		jsonError(w, "internal server error", http.StatusInternalServerError)
		return
	}

	if !owner.IsActive {
		jsonError(w, "account is inactive", http.StatusForbidden)
		return
	}

	if err = bcrypt.CompareHashAndPassword([]byte(owner.PasswordHash), []byte(req.Password)); err != nil {
		jsonError(w, "invalid username or password", http.StatusUnauthorized)
		return
	}

	_ = h.repo.UpdateLastLogin(r.Context(), owner.ID)

	token, err := h.store.CreateToken(domain.OwnerSessionData{OwnerID: owner.ID})
	if err != nil {
		jsonError(w, "could not create session", http.StatusInternalServerError)
		return
	}

	jsonResponse(w, http.StatusOK, map[string]any{
		"message": "login successful",
		"token":   token,
	})
}

// POST /api/v1/owner/auth/logout
func (h *OwnerHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// JWT is stateless — client simply discards the token.
	jsonResponse(w, http.StatusOK, map[string]any{"message": "logged out"})
}

// GET /api/v1/owner/auth/me
func (h *OwnerHandler) Me(w http.ResponseWriter, r *http.Request) {
	sd := middleware.GetOwnerSession(r)
	if sd == nil {
		jsonError(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	owner, err := h.repo.FindByID(r.Context(), sd.OwnerID)
	if err != nil {
		jsonError(w, "owner not found", http.StatusNotFound)
		return
	}

	jsonResponse(w, http.StatusOK, domain.OwnerMeResponse{
		ID:       owner.ID,
		FullName: owner.FullName,
		Email:    owner.Email,
		Username: owner.Username,
		Role:     "owner",
	})
}

// GET /api/v1/owner/tenants
func (h *OwnerHandler) ListTenants(w http.ResponseWriter, r *http.Request) {
	tenants, err := h.repo.ListTenants(r.Context())
	if err != nil {
		jsonError(w, "failed to list tenants", http.StatusInternalServerError)
		return
	}
	if tenants == nil {
		tenants = []domain.Tenant{}
	}
	jsonResponse(w, http.StatusOK, map[string]any{"data": tenants, "total": len(tenants)})
}

// GET /api/v1/owner/stats
func (h *OwnerHandler) Stats(w http.ResponseWriter, r *http.Request) {
	stats, err := h.repo.PlatformStats(r.Context())
	if err != nil {
		jsonError(w, "failed to fetch stats", http.StatusInternalServerError)
		return
	}
	jsonResponse(w, http.StatusOK, stats)
}
