package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/linkchat/api/internal/domain"
	"github.com/linkchat/api/internal/middleware"
	"github.com/linkchat/api/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	repo     *repository.UserRepository
	authRepo *repository.AuthRepository
}

func NewUserHandler(repo *repository.UserRepository, authRepo *repository.AuthRepository) *UserHandler {
	return &UserHandler{repo: repo, authRepo: authRepo}
}

func (h *UserHandler) List(w http.ResponseWriter, r *http.Request) {
	sd := middleware.GetSession(r)
	users, err := h.repo.List(r.Context(), sd.TenantID)
	if err != nil {
		jsonError(w, "failed to list users", http.StatusInternalServerError)
		return
	}
	// Strip password hashes from response
	type safeUser struct {
		ID        uint64 `json:"id"`
		FullName  string `json:"full_name"`
		Email     string `json:"email"`
		Username  string `json:"username"`
		Role      string `json:"role"`
		IsActive  bool   `json:"is_active"`
	}
	safe := make([]safeUser, 0, len(users))
	for _, u := range users {
		safe = append(safe, safeUser{u.ID, u.FullName, u.Email, u.Username, u.Role, u.IsActive})
	}
	jsonResponse(w, http.StatusOK, map[string]any{"data": safe})
}

func (h *UserHandler) Invite(w http.ResponseWriter, r *http.Request) {
	sd := middleware.GetSession(r)
	if sd.Role != "admin" {
		jsonError(w, "only admins can invite users", http.StatusForbidden)
		return
	}
	var req domain.InviteUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "invalid body", http.StatusBadRequest)
		return
	}
	req.Username = strings.ToLower(strings.TrimSpace(req.Username))
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))

	if req.Username == "" || req.Password == "" || req.Role == "" {
		jsonError(w, "username, password, and role are required", http.StatusUnprocessableEntity)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
	if err != nil {
		jsonError(w, "internal error", http.StatusInternalServerError)
		return
	}
	user, err := h.repo.Create(r.Context(), domain.TenantUser{
		TenantID: sd.TenantID, FullName: req.FullName, Email: req.Email,
		Username: req.Username, PasswordHash: string(hash), Role: req.Role,
	})
	if err != nil {
		jsonError(w, "failed to invite user", http.StatusInternalServerError)
		return
	}
	jsonResponse(w, http.StatusCreated, map[string]any{
		"id": user.ID, "username": user.Username, "role": user.Role,
	})
}

func (h *UserHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	sd := middleware.GetSession(r)
	var req domain.UpdateUserRequest
	json.NewDecoder(r.Body).Decode(&req)
	if err := h.repo.UpdateProfile(r.Context(), sd.UserID, sd.TenantID, req.FullName, req.Email); err != nil {
		jsonError(w, "failed to update profile", http.StatusInternalServerError)
		return
	}
	jsonResponse(w, http.StatusOK, map[string]string{"message": "profile updated"})
}

func (h *UserHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	sd := middleware.GetSession(r)
	var req domain.ChangePasswordRequest
	json.NewDecoder(r.Body).Decode(&req)

	user, err := h.repo.FindByID(r.Context(), sd.UserID, sd.TenantID)
	if err != nil {
		jsonError(w, "user not found", http.StatusNotFound)
		return
	}
	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.CurrentPassword)) != nil {
		jsonError(w, "current password is incorrect", http.StatusUnauthorized)
		return
	}
	hash, _ := bcrypt.GenerateFromPassword([]byte(req.NewPassword), 12)
	h.repo.UpdatePassword(r.Context(), sd.UserID, string(hash))
	jsonResponse(w, http.StatusOK, map[string]string{"message": "password changed"})
}

func (h *UserHandler) ToggleActive(w http.ResponseWriter, r *http.Request) {
	sd := middleware.GetSession(r)
	if sd.Role != "admin" {
		jsonError(w, "forbidden", http.StatusForbidden)
		return
	}
	id, _ := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	var body struct{ Active bool `json:"is_active"` }
	json.NewDecoder(r.Body).Decode(&body)
	h.repo.ToggleActive(r.Context(), id, sd.TenantID, body.Active)
	jsonResponse(w, http.StatusOK, map[string]string{"message": "updated"})
}

func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	sd := middleware.GetSession(r)
	if sd.Role != "admin" {
		jsonError(w, "forbidden", http.StatusForbidden)
		return
	}
	id, _ := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if id == sd.UserID {
		jsonError(w, "cannot delete yourself", http.StatusBadRequest)
		return
	}
	if err := h.repo.Delete(r.Context(), id, sd.TenantID); errors.Is(err, repository.ErrNotFound) {
		jsonError(w, "user not found", http.StatusNotFound)
		return
	}
	jsonResponse(w, http.StatusOK, map[string]string{"message": "deleted"})
}
