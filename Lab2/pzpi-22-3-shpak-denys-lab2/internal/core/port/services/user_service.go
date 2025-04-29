package services // import "wayra/internal/core/port/services"

import (
	"wayra/internal/core/domain/models"
)

// UserService is a service that manages the user domain model
type UserService interface {
	Service[models.User]
}
