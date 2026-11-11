package middleware

import (
	"context"
	"net/http"

	"github.com/linkchat/api/internal/domain"
	"github.com/linkchat/api/internal/session"
)

type ownerContextKey string

const OwnerSessionKey ownerContextKey = "owner_session"

// RequireOwner is an HTTP middleware that validates the owner JWT.
// Completely separate from RequireAuth (tenant middleware).
func RequireOwner(store *session.OwnerStore) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			sd := store.Get(r)
			if sd == nil {
				http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), OwnerSessionKey, sd)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetOwnerSession retrieves the owner session from request context.
func GetOwnerSession(r *http.Request) *domain.OwnerSessionData {
	sd, _ := r.Context().Value(OwnerSessionKey).(*domain.OwnerSessionData)
	return sd
}
