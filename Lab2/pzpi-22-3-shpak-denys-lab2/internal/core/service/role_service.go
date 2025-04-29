package service // import "wayra/internal/core/service"

import (
	"wayra/internal/core/domain/models"
	"wayra/internal/core/port"
	"wayra/internal/core/port/services"
)

// RoleService is a struct to handle the Role service
type RoleService struct {
	*GenericService[models.Role] // Embedding the GenericService
}

// NewRoleService is a function to create a new RoleService
// repo: Repository of Role
// returns: RoleService
func NewRoleService(repo port.Repository[models.Role]) services.RoleService {
	return &RoleService{
		GenericService: NewGenericService(repo),
	}
}
