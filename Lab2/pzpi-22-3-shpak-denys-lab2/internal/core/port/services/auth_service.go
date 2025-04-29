// Package services provides the service interfaces for the core domain.
package services // import "wayra/internal/core/port/services"

import (
	"context"
	"wayra/internal/core/domain/models"

	"github.com/golang-jwt/jwt/v5"
)

// AuthService provides the service interface for the authentication service.
type AuthService interface {
	RegisterUser(ctx context.Context, user *models.User) error
	LoginUser(ctx context.Context, username, password string) (string, error)
	ValidateToken(token string) (*jwt.RegisteredClaims, error)
}
