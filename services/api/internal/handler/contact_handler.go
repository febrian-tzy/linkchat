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

type ContactHandler struct{ repo *repository.ContactRepository }

func NewContactHandler(repo *repository.ContactRepository) *ContactHandler {
	return &ContactHandler{repo: repo}
}

func (h *ContactHandler) List(w http.ResponseWriter, r *http.Request) {
	sd := middleware.GetSession(r)
	search := r.URL.Query().Get("search")
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	if limit <= 0 { limit = 20 }

	contacts, total, err := h.repo.List(r.Context(), sd.TenantID, search, limit, offset)
	if err != nil {
		jsonError(w, "failed to list contacts", http.StatusInternalServerError)
		return
	}
	if contacts == nil { contacts = []domain.Contact{} }
	jsonResponse(w, http.StatusOK, map[string]any{"data": contacts, "total": total})
}

func (h *ContactHandler) Get(w http.ResponseWriter, r *http.Request) {
	sd := middleware.GetSession(r)
	id, _ := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	c, err := h.repo.FindByID(r.Context(), id, sd.TenantID)
	if errors.Is(err, repository.ErrNotFound) {
		jsonError(w, "contact not found", http.StatusNotFound)
		return
	}
	jsonResponse(w, http.StatusOK, c)
}

func (h *ContactHandler) Create(w http.ResponseWriter, r *http.Request) {
	sd := middleware.GetSession(r)
	var req domain.CreateContactRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "invalid body", http.StatusBadRequest)
		return
	}
	if req.Name == "" || req.Phone == "" {
		jsonError(w, "name and phone are required", http.StatusUnprocessableEntity)
		return
	}
	c, err := h.repo.Create(r.Context(), domain.Contact{
		TenantID: sd.TenantID, Name: req.Name, Phone: req.Phone, Email: req.Email, Notes: req.Notes,
	})
	if err != nil {
		jsonError(w, "failed to create contact", http.StatusInternalServerError)
		return
	}
	jsonResponse(w, http.StatusCreated, c)
}

func (h *ContactHandler) Update(w http.ResponseWriter, r *http.Request) {
	sd := middleware.GetSession(r)
	id, _ := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	var req domain.UpdateContactRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "invalid body", http.StatusBadRequest)
		return
	}
	c, err := h.repo.Update(r.Context(), domain.Contact{
		ID: id, TenantID: sd.TenantID, Name: req.Name, Phone: req.Phone, Email: req.Email, Notes: req.Notes, IsBlocked: req.IsBlocked,
	})
	if errors.Is(err, repository.ErrNotFound) {
		jsonError(w, "contact not found", http.StatusNotFound)
		return
	}
	if err != nil {
		jsonError(w, "failed to update contact", http.StatusInternalServerError)
		return
	}
	jsonResponse(w, http.StatusOK, c)
}

func (h *ContactHandler) Delete(w http.ResponseWriter, r *http.Request) {
	sd := middleware.GetSession(r)
	id, _ := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err := h.repo.Delete(r.Context(), id, sd.TenantID); errors.Is(err, repository.ErrNotFound) {
		jsonError(w, "contact not found", http.StatusNotFound)
		return
	}
	jsonResponse(w, http.StatusOK, map[string]string{"message": "deleted"})
}
