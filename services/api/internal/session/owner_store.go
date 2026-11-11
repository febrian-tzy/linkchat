package session

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/linkchat/api/internal/domain"
)

type ownerClaims struct {
	OwnerID int64 `json:"oid"` // int64 to avoid JSON float64 parsing issues
	jwt.RegisteredClaims
}

// OwnerStore handles JWT-based auth for platform owners.
type OwnerStore struct {
	secret []byte
}

func NewOwnerStore(secret string) *OwnerStore {
	combined := "owner_" + secret
	fmt.Printf("[OwnerStore] initialized with secret length: %d\n", len(combined))
	return &OwnerStore{secret: []byte(combined)}
}

// CreateToken signs a JWT for the given owner session.
func (s *OwnerStore) CreateToken(data domain.OwnerSessionData) (string, error) {
	now := time.Now()
	c := ownerClaims{
		OwnerID: int64(data.OwnerID),
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(12 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	signed, err := token.SignedString(s.secret)
	if err != nil {
		return "", fmt.Errorf("ownerStore.CreateToken: %w", err)
	}
	return signed, nil
}

// Get reads and validates a Bearer token from the Authorization header.
func (s *OwnerStore) Get(r *http.Request) *domain.OwnerSessionData {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		fmt.Println("[OwnerStore] no Authorization header")
		return nil
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
		fmt.Println("[OwnerStore] bad Authorization format")
		return nil
	}

	fmt.Printf("[OwnerStore] parsing token: %s...\n", parts[1][:20])

	t, err := jwt.ParseWithClaims(parts[1], &ownerClaims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return s.secret, nil
	})
	if err != nil {
		fmt.Printf("[OwnerStore] parse error: %v\n", err)
		return nil
	}
	if !t.Valid {
		fmt.Println("[OwnerStore] token not valid")
		return nil
	}

	c, ok := t.Claims.(*ownerClaims)
	if !ok {
		fmt.Println("[OwnerStore] claims type assertion failed")
		return nil
	}
	fmt.Printf("[OwnerStore] OwnerID from claims: %d\n", c.OwnerID)
	if c.OwnerID <= 0 {
		fmt.Println("[OwnerStore] OwnerID is zero")
		return nil
	}

	return &domain.OwnerSessionData{OwnerID: uint64(c.OwnerID)}
}
