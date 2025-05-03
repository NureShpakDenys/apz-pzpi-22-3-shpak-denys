// Package handlers contains the http handlers for the application
//
// This package is part of the adapter layer
// You can get handler instances by calling the New*Handler functions
// The handlers are responsible for handling the http requests and responses
package handlers // import "wayra/internal/adapter/httpserver/handlers"
import (
	"errors"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// The role of user is one of these constants
const (
	AdminRole = iota + 1
	UserRole
	ManagerRole
)

// Role is a type for the user role in company
type Role string

// The role of user in company is one of these constants
const (
	RoleUser        Role = "user"
	RoleAdmin       Role = "admin"
	RoleManager     Role = "manager"
	RoleDBAdmin     Role = "db_admin"
	RoleSystemAdmin Role = "system_admin"
)

// The status of delivery is one of these constants
const (
	NotStarted = "not_started"
	InProgress = "in_progress"
	Completed  = "completed"
)

// GetUserIDFromToken gets the user ID from the token in the request
// c: The gin context
// Returns: The user ID and an error if there was a problem
// GetUserIDFromToken gets the user ID from the token in the request
// c: The gin context
// Returns: The user ID and an error if there was a problem
func GetUserIDFromToken(c *gin.Context) (*uint, error) {
	tokenCookie, exists := c.Get("token")
	if !exists {
		return nil, errors.New("token not found in context")
	}

	tokenString, ok := tokenCookie.(string)
	if !ok {
		return nil, errors.New("invalid token format in context")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}

		return []byte("mysecret123"), nil
	})
	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid token: %s", err.Error())
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		stringUserID, ok := claims["sub"].(string)
		if !ok {
			return nil, errors.New("sub not found in token claims")
		}

		userID, err := strconv.Atoi(stringUserID)
		if err != nil {
			return nil, errors.New("problem parsing user ID")
		}

		uintUserID := uint(userID)
		return &uintUserID, nil
	}

	return nil, errors.New("failed to parse token claims")
}
