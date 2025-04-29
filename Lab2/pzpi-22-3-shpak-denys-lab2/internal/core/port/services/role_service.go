package services // import "wayra/internal/core/port/services"

import (
	"wayra/internal/core/domain/models"
)

// RoleService is an interface to represent the role service
type RoleService interface {
	Service[models.Role]
}
