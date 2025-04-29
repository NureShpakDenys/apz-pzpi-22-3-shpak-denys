// Package service implements the core business logic of the application.
package service // import "wayra/internal/core/service"

import (
	"context"
	"errors"
	"strconv"
	"time"

	"wayra/internal/core/domain/models"
	"wayra/internal/core/port/services"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// AuthService is a service that provides authentication and authorization functionality.
type AuthService struct {
	userService services.UserService // User service to interact with the user repository.
	secretKey   string               // Secret key used to sign the JWT tokens.
	tokenExpiry time.Duration        // Token expiry time.
}

// NewAuthService creates a new instance of the AuthService.
// userService: User service to interact with the user repository.
// secretKey: Secret key used to sign the JWT tokens.
// tokenExpiry: Token expiry time.
// returns: A new instance of the AuthService.
func NewAuthService(userService services.UserService, secretKey string, tokenExpiry time.Duration) services.AuthService {
	return &AuthService{
		userService: userService,
		secretKey:   secretKey,
		tokenExpiry: tokenExpiry,
	}
}

// CustomClaims represents the custom claims of the JWT token.
type CustomClaims struct {
	jwt.RegisteredClaims        // Standard claims.
	Username             string `json:"username"` // Username of the user.
}

// RegisterUser registers a new user.
// ctx: Context of the request.
// user: User to register.
// returns: An error if the operation failed.
func (s *AuthService) RegisterUser(ctx context.Context, user *models.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)

	return s.userService.Create(ctx, user)
}

// LoginUser logs in a user.
// ctx: Context of the request.
// username: Username of the user.
// password: Password of the user.
// returns: A JWT token if the operation was successful.
func (s *AuthService) LoginUser(ctx context.Context, username, password string) (string, error) {
	users, err := s.userService.Where(ctx, &models.User{Name: username})
	if err != nil || len(users) == 0 {
		return "", errors.New("invalid username or password")
	}

	user := users[0]
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("invalid username or password")
	}

	token, err := s.generateToken(user.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}

// ValidateToken validates a JWT token.
// token: JWT token to validate.
// returns: The claims of the token if the operation was successful.
func (s *AuthService) ValidateToken(token string) (*jwt.RegisteredClaims, error) {
	parsedToken, err := jwt.ParseWithClaims(token, &CustomClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(s.secretKey), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := parsedToken.Claims.(*CustomClaims)
	if !ok || !parsedToken.Valid {
		return nil, errors.New("invalid token")
	}

	return &claims.RegisteredClaims, nil
}

// generateToken generates a JWT token.
// userID: ID of the user.
// returns: A JWT token if the operation was successful.
func (s *AuthService) generateToken(userID uint) (string, error) {
	claims := CustomClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   strconv.Itoa(int(userID)),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.tokenExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        strconv.Itoa(int(userID)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.secretKey))
}
