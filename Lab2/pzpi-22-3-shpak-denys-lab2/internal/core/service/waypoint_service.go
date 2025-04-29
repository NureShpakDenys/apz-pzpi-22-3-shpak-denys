package service // import "wayra/internal/core/service"

import (
	"wayra/internal/core/domain/models"
	"wayra/internal/core/port"
	"wayra/internal/core/port/services"
)

// WaypointService is a service that manages waypoints
type WaypointService struct {
	*GenericService[models.Waypoint] // Embedding the generic service
}

// NewWaypointService creates a new waypoint service
// repo: the repository to use
// returns: a new waypoint service
func NewWaypointService(repo port.Repository[models.Waypoint]) services.WaypointService {
	return &WaypointService{
		GenericService: NewGenericService(repo),
	}
}
