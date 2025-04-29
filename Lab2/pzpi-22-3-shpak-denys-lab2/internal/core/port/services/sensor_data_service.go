package services // import "wayra/internal/core/port/services"

import (
	"wayra/internal/core/domain/models"
)

// SensorDataService is the interface that wraps the basic SensorData service methods.
type SensorDataService interface {
	Service[models.SensorData]
}
