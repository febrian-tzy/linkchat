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

type ConversationHandler struct{ repo *repository.ConversationRepository }

func NewConversationHandler(repo *repository.ConversationRepository) *ConversationHandler {
	return &ConversationHandler{repo: repo}
}

func (h *ConversationHandler) List(w http.ResponseWriter, r *http.Request) {
	sd := middleware.GetSession(r)
	status := r.URL.Query().Get("status")
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	if limit <= 0 { limit = 20 }

	convs, total, err := h.repo.List(r.Context(), sd.TenantID, status, limit, offset)
	if err != nil {
		jsonError(w, "failed to list conversations", http.StatusInternalServerError)
		return
	}
	if convs == nil { convs = []domain.Conversation{} }
	jsonResponse(w, http.StatusOK, map[string]any{"data": convs, "total": total})
}

func (h *ConversationHandler) Get(w http.ResponseWriter, r *http.Request) {
	sd := middleware.GetSession(r)
	id, _ := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	cv, err := h.repo.FindByID(r.Context(), id, sd.TenantID)
	if errors.Is(err, repository.ErrNotFound) {
		jsonError(w, "conversation not found", http.StatusNotFound)
		return
	}
	jsonResponse(w, http.StatusOK, cv)
}

func (h *ConversationHandler) Create(w http.ResponseWriter, r *http.Request) {
	sd := middleware.GetSession(r)
	var req domain.CreateConversationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "invalid body", http.StatusBadRequest)
		return
	}
	cv, err := h.repo.Create(r.Context(), domain.Conversation{
		TenantID: sd.TenantID, InboxID: req.InboxID, ContactID: req.ContactID,
	})
	if err != nil {
		jsonError(w, "failed to create conversation", http.StatusInternalServerError)
		return
	}
	jsonResponse(w, http.StatusCreated, cv)
}

func (h *ConversationHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	sd := middleware.GetSession(r)
	id, _ := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	var req domain.UpdateConversationRequest
	json.NewDecoder(r.Body).Decode(&req)
	if req.Status == "" { req.Status = "open" }
	if err := h.repo.UpdateStatus(r.Context(), id, sd.TenantID, req.Status, req.AssignedUserID); err != nil {
		jsonError(w, "failed to update", http.StatusInternalServerError)
		return
	}
	jsonResponse(w, http.StatusOK, map[string]string{"message": "updated"})
}

func (h *ConversationHandler) ListMessages(w http.ResponseWriter, r *http.Request) {
	sd := middleware.GetSession(r)
	id, _ := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	if limit <= 0 { limit = 50 }

	msgs, total, err := h.repo.ListMessages(r.Context(), id, sd.TenantID, limit, offset)
	if err != nil {
		jsonError(w, "failed to list messages", http.StatusInternalServerError)
		return
	}
	if msgs == nil { msgs = []domain.Message{} }
	jsonResponse(w, http.StatusOK, map[string]any{"data": msgs, "total": total})
}

func (h *ConversationHandler) SendMessage(w http.ResponseWriter, r *http.Request) {
	sd := middleware.GetSession(r)
	convID, _ := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	var req domain.SendMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "invalid body", http.StatusBadRequest)
		return
	}
	if req.Content == "" {
		jsonError(w, "content is required", http.StatusUnprocessableEntity)
		return
	}
	if req.ContentType == "" { req.ContentType = "text" }
	msg, err := h.repo.CreateMessage(r.Context(), domain.Message{
		ConversationID: convID, TenantID: sd.TenantID,
		SenderType: "agent", SenderID: &sd.UserID,
		Content: req.Content, ContentType: req.ContentType,
	})
	if err != nil {
		jsonError(w, "failed to send message", http.StatusInternalServerError)
		return
	}
	jsonResponse(w, http.StatusCreated, msg)
}
