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

type AutomationHandler struct{ repo *repository.AutomationRepository }

func NewAutomationHandler(repo *repository.AutomationRepository) *AutomationHandler {
	return &AutomationHandler{repo: repo}
}

func (h *AutomationHandler) List(w http.ResponseWriter, r *http.Request) {
	sd := middleware.GetSession(r)
	list, err := h.repo.List(r.Context(), sd.TenantID)
	if err != nil {
		jsonError(w, "failed to list automations", http.StatusInternalServerError)
		return
	}
	if list == nil { list = []domain.Automation{} }
	jsonResponse(w, http.StatusOK, map[string]any{"data": list})
}

func (h *AutomationHandler) Get(w http.ResponseWriter, r *http.Request) {
	sd := middleware.GetSession(r)
	id, _ := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	a, err := h.repo.FindByID(r.Context(), id, sd.TenantID)
	if errors.Is(err, repository.ErrNotFound) {
		jsonError(w, "automation not found", http.StatusNotFound)
		return
	}
	jsonResponse(w, http.StatusOK, a)
}

func (h *AutomationHandler) Create(w http.ResponseWriter, r *http.Request) {
	sd := middleware.GetSession(r)
	var req domain.CreateAutomationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "invalid body", http.StatusBadRequest)
		return
	}
	if req.Name == "" || req.Event == "" {
		jsonError(w, "name and event are required", http.StatusUnprocessableEntity)
		return
	}
	a, err := h.repo.Create(r.Context(), domain.Automation{
		TenantID: sd.TenantID, Name: req.Name, Description: req.Description,
		Event: req.Event, Conditions: req.Conditions, Actions: req.Actions,
	})
	if err != nil {
		jsonError(w, "failed to create automation", http.StatusInternalServerError)
		return
	}
	jsonResponse(w, http.StatusCreated, a)
}

func (h *AutomationHandler) Toggle(w http.ResponseWriter, r *http.Request) {
	sd := middleware.GetSession(r)
	id, _ := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	var body struct{ Active bool `json:"is_active"` }
	json.NewDecoder(r.Body).Decode(&body)
	if err := h.repo.ToggleActive(r.Context(), id, sd.TenantID, body.Active); err != nil {
		jsonError(w, "failed to update automation", http.StatusInternalServerError)
		return
	}
	jsonResponse(w, http.StatusOK, map[string]string{"message": "updated"})
}

func (h *AutomationHandler) Delete(w http.ResponseWriter, r *http.Request) {
	sd := middleware.GetSession(r)
	id, _ := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err := h.repo.Delete(r.Context(), id, sd.TenantID); errors.Is(err, repository.ErrNotFound) {
		jsonError(w, "automation not found", http.StatusNotFound)
		return
	}
	jsonResponse(w, http.StatusOK, map[string]string{"message": "deleted"})
}
