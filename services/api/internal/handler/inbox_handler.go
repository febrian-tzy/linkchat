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

type InboxHandler struct{ repo *repository.InboxRepository }

func NewInboxHandler(repo *repository.InboxRepository) *InboxHandler {
	return &InboxHandler{repo: repo}
}

func (h *InboxHandler) List(w http.ResponseWriter, r *http.Request) {
	sd := middleware.GetSession(r)
	inboxes, err := h.repo.List(r.Context(), sd.TenantID)
	if err != nil {
		jsonError(w, "failed to list inboxes", http.StatusInternalServerError)
		return
	}
	if inboxes == nil { inboxes = []domain.Inbox{} }
	jsonResponse(w, http.StatusOK, map[string]any{"data": inboxes})
}

func (h *InboxHandler) Get(w http.ResponseWriter, r *http.Request) {
	sd := middleware.GetSession(r)
	id, _ := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	inbox, err := h.repo.FindByID(r.Context(), id, sd.TenantID)
	if errors.Is(err, repository.ErrNotFound) {
		jsonError(w, "inbox not found", http.StatusNotFound)
		return
	}
	jsonResponse(w, http.StatusOK, inbox)
}

func (h *InboxHandler) Create(w http.ResponseWriter, r *http.Request) {
	sd := middleware.GetSession(r)
	var req domain.CreateInboxRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "invalid body", http.StatusBadRequest)
		return
	}
	if req.Name == "" {
		jsonError(w, "name is required", http.StatusUnprocessableEntity)
		return
	}
	if req.ChannelType == "" { req.ChannelType = "whatsapp" }
	inbox, err := h.repo.Create(r.Context(), domain.Inbox{
		TenantID: sd.TenantID, Name: req.Name, ChannelType: req.ChannelType,
	})
	if err != nil {
		jsonError(w, "failed to create inbox", http.StatusInternalServerError)
		return
	}
	jsonResponse(w, http.StatusCreated, inbox)
}

func (h *InboxHandler) Delete(w http.ResponseWriter, r *http.Request) {
	sd := middleware.GetSession(r)
	id, _ := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err := h.repo.Delete(r.Context(), id, sd.TenantID); errors.Is(err, repository.ErrNotFound) {
		jsonError(w, "inbox not found", http.StatusNotFound)
		return
	}
	jsonResponse(w, http.StatusOK, map[string]string{"message": "deleted"})
}
