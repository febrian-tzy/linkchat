package handler

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/linkchat/api/internal/domain"
	"github.com/linkchat/api/internal/middleware"
	"github.com/linkchat/api/internal/repository"
)

type WebhookHandler struct{ repo *repository.WebhookRepository }
type APIKeyHandler struct{ repo *repository.APIKeyRepository }

func NewWebhookHandler(repo *repository.WebhookRepository) *WebhookHandler {
	return &WebhookHandler{repo: repo}
}
func NewAPIKeyHandler(repo *repository.APIKeyRepository) *APIKeyHandler {
	return &APIKeyHandler{repo: repo}
}

// ── Webhooks ─────────────────────────────────────────────────

func (h *WebhookHandler) List(w http.ResponseWriter, r *http.Request) {
	sd := middleware.GetSession(r)
	list, err := h.repo.List(r.Context(), sd.TenantID)
	if err != nil {
		jsonError(w, "failed to list webhooks", http.StatusInternalServerError)
		return
	}
	if list == nil { list = []domain.Webhook{} }
	jsonResponse(w, http.StatusOK, map[string]any{"data": list})
}

func (h *WebhookHandler) Create(w http.ResponseWriter, r *http.Request) {
	sd := middleware.GetSession(r)
	var req domain.CreateWebhookRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "invalid body", http.StatusBadRequest)
		return
	}
	if req.URL == "" {
		jsonError(w, "url is required", http.StatusUnprocessableEntity)
		return
	}
	wh, err := h.repo.Create(r.Context(), domain.Webhook{
		TenantID: sd.TenantID, URL: req.URL, Secret: req.Secret, Events: req.Events,
	})
	if err != nil {
		jsonError(w, "failed to create webhook", http.StatusInternalServerError)
		return
	}
	jsonResponse(w, http.StatusCreated, wh)
}

func (h *WebhookHandler) Delete(w http.ResponseWriter, r *http.Request) {
	sd := middleware.GetSession(r)
	id, _ := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err := h.repo.Delete(r.Context(), id, sd.TenantID); errors.Is(err, repository.ErrNotFound) {
		jsonError(w, "webhook not found", http.StatusNotFound)
		return
	}
	jsonResponse(w, http.StatusOK, map[string]string{"message": "deleted"})
}

// ── API Keys ─────────────────────────────────────────────────

func (h *APIKeyHandler) List(w http.ResponseWriter, r *http.Request) {
	sd := middleware.GetSession(r)
	list, err := h.repo.List(r.Context(), sd.TenantID)
	if err != nil {
		jsonError(w, "failed to list api keys", http.StatusInternalServerError)
		return
	}
	if list == nil { list = []domain.APIKey{} }
	jsonResponse(w, http.StatusOK, map[string]any{"data": list})
}

func (h *APIKeyHandler) Create(w http.ResponseWriter, r *http.Request) {
	sd := middleware.GetSession(r)
	var req domain.CreateAPIKeyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "invalid body", http.StatusBadRequest)
		return
	}
	if req.Name == "" {
		jsonError(w, "name is required", http.StatusUnprocessableEntity)
		return
	}

	// Generate a random 32-byte key
	raw := make([]byte, 32)
	rand.Read(raw)
	plainKey := "lc_" + hex.EncodeToString(raw)

	hash := sha256.Sum256([]byte(plainKey))
	keyHash := hex.EncodeToString(hash[:])
	prefix := plainKey[:10]

	key, err := h.repo.Create(r.Context(), domain.APIKey{
		TenantID:  sd.TenantID,
		Name:      req.Name,
		KeyHash:   keyHash,
		KeyPrefix: prefix,
		IsActive:  true,
	})
	if err != nil {
		jsonError(w, "failed to create api key", http.StatusInternalServerError)
		return
	}

	jsonResponse(w, http.StatusCreated, map[string]any{
		"id":         key.ID,
		"name":       key.Name,
		"key_prefix": key.KeyPrefix,
		"key":        plainKey, // shown only once
		"created_at": key.CreatedAt,
		"message":    fmt.Sprintf("Save this key now. It will not be shown again: %s", plainKey),
	})
}

func (h *APIKeyHandler) Delete(w http.ResponseWriter, r *http.Request) {
	sd := middleware.GetSession(r)
	id, _ := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err := h.repo.Delete(r.Context(), id, sd.TenantID); errors.Is(err, repository.ErrNotFound) {
		jsonError(w, "api key not found", http.StatusNotFound)
		return
	}
	jsonResponse(w, http.StatusOK, map[string]string{"message": "deleted"})
}
