package middleware

import (
	"context"
	"net/http"

	"github.com/linkchat/api/internal/domain"
	"github.com/linkchat/api/internal/session"
)

type contextKey string

// SessionKey is the context key for the authenticated session data.
const SessionKey contextKey = "session"

// RequireAuth is an HTTP middleware that checks for a valid server-side session.
// If the session is missing or invalid, it returns 401 Unauthorized.
func RequireAuth(store *session.Store) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			sd := store.Get(r)
			if sd == nil {
				http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
				return
			}
			// Inject session data into request context.
			ctx := context.WithValue(r.Context(), SessionKey, sd)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetSession retrieves the session data from the request context.
func GetSession(r *http.Request) *domain.SessionData {
	sd, _ := r.Context().Value(SessionKey).(*domain.SessionData)
	return sd
}
