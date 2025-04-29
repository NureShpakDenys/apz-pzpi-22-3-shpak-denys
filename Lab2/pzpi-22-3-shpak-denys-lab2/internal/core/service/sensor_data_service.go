package service // import "wayra/internal/core/service"

import (
	"wayra/internal/core/domain/models"
	"wayra/internal/core/port"
	"wayra/internal/core/port/services"
)

// SensorDataService is a service that manages the sensor data
type SensorDataService struct {
	*GenericService[models.SensorData] // Embedding the generic service
}

// NewSensorDataService creates a new sensor data service
// repo: the repository to use
// returns: a new sensor data service
func NewSensorDataService(repo port.Repository[models.SensorData]) services.SensorDataService {
	return &SensorDataService{
		GenericService: NewGenericService(repo),
	}
}
