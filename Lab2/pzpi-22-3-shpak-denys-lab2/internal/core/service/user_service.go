package service // import "wayra/internal/core/service"

import (
	"wayra/internal/core/domain/models"
	"wayra/internal/core/port"
	"wayra/internal/core/port/services"
)

// UserService is a struct that defines the service for the user entity
type UserService struct {
	*GenericService[models.User] // Embedding the generic service
}

// NewUserService is a function that creates a new user service
// repo: Repository for the user entity
// Returns: A new user service
func NewUserService(repo port.Repository[models.User]) services.UserService {
	return &UserService{
		GenericService: NewGenericService(repo),
	}
}
