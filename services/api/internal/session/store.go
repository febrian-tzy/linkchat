package session

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/linkchat/api/internal/domain"
)

type claims struct {
	UserID   int64  `json:"uid"` // int64 avoids JSON float64 parsing issues
	TenantID int64  `json:"tid"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// Store handles JWT-based auth via Authorization header.
type Store struct {
	secret []byte
}

// New creates a new Store.
func New(secret string) *Store {
	return &Store{secret: []byte(secret)}
}

// CreateToken generates a signed JWT string for the given session data.
func (s *Store) CreateToken(data domain.SessionData) (string, error) {
	now := time.Now()
	c := claims{
		UserID:   int64(data.UserID),
		TenantID: int64(data.TenantID),
		Role:     data.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(7 * 24 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	signed, err := token.SignedString(s.secret)
	if err != nil {
		return "", fmt.Errorf("session.CreateToken: %w", err)
	}
	return signed, nil
}

// Save is kept for interface compatibility.
func (s *Store) Save(w http.ResponseWriter, r *http.Request, data domain.SessionData) error {
	_, err := s.CreateToken(data)
	return err
}

// Get reads JWT from Authorization: Bearer <token> header.
func (s *Store) Get(r *http.Request) *domain.SessionData {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return nil
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
		return nil
	}

	t, err := jwt.ParseWithClaims(parts[1], &claims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return s.secret, nil
	})
	if err != nil || !t.Valid {
		return nil
	}

	c, ok := t.Claims.(*claims)
	if !ok || c.UserID <= 0 || c.TenantID <= 0 {
		return nil
	}

	return &domain.SessionData{
		UserID:   uint64(c.UserID),
		TenantID: uint64(c.TenantID),
		Role:     c.Role,
	}
}

// Clear is a no-op for JWT — client deletes token from localStorage.
func (s *Store) Clear(w http.ResponseWriter, r *http.Request) error {
	return nil
}
