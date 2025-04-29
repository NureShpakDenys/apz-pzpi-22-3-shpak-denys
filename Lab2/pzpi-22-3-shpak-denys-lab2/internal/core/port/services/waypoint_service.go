package services // import "wayra/internal/core/port/services"

import (
	"wayra/internal/core/domain/models"
)

// WaypointService is the interface that wraps the basic Waypoint methods.
type WaypointService interface {
	Service[models.Waypoint]
}
