package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/linkchat/api/internal/domain"
	"github.com/linkchat/api/internal/middleware"
	"github.com/linkchat/api/internal/repository"
)

type BroadcastHandler struct{ repo *repository.BroadcastRepository }

func NewBroadcastHandler(repo *repository.BroadcastRepository) *BroadcastHandler {
	return &BroadcastHandler{repo: repo}
}

func (h *BroadcastHandler) List(w http.ResponseWriter, r *http.Request) {
	sd := middleware.GetSession(r)
	list, err := h.repo.List(r.Context(), sd.TenantID)
	if err != nil {
		jsonError(w, "failed to list broadcasts", http.StatusInternalServerError)
		return
	}
	if list == nil { list = []domain.Broadcast{} }
	jsonResponse(w, http.StatusOK, map[string]any{"data": list})
}

func (h *BroadcastHandler) Get(w http.ResponseWriter, r *http.Request) {
	sd := middleware.GetSession(r)
	id, _ := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	b, err := h.repo.FindByID(r.Context(), id, sd.TenantID)
	if errors.Is(err, repository.ErrNotFound) {
		jsonError(w, "broadcast not found", http.StatusNotFound)
		return
	}
	jsonResponse(w, http.StatusOK, b)
}

func (h *BroadcastHandler) Create(w http.ResponseWriter, r *http.Request) {
	sd := middleware.GetSession(r)
	var req domain.CreateBroadcastRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "invalid body", http.StatusBadRequest)
		return
	}
	if req.Name == "" || req.Message == "" {
		jsonError(w, "name and message are required", http.StatusUnprocessableEntity)
		return
	}
	b, err := h.repo.Create(r.Context(), domain.Broadcast{
		TenantID: sd.TenantID, InboxID: req.InboxID,
		Name: req.Name, Message: req.Message, CreatedBy: sd.UserID,
	}, req.ContactIDs)
	if err != nil {
		jsonError(w, "failed to create broadcast", http.StatusInternalServerError)
		return
	}
	jsonResponse(w, http.StatusCreated, b)
}

func (h *BroadcastHandler) Delete(w http.ResponseWriter, r *http.Request) {
	sd := middleware.GetSession(r)
	id, _ := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err := h.repo.Delete(r.Context(), id, sd.TenantID); errors.Is(err, repository.ErrNotFound) {
		jsonError(w, "broadcast not found", http.StatusNotFound)
		return
	}
	jsonResponse(w, http.StatusOK, map[string]string{"message": "deleted"})
}
