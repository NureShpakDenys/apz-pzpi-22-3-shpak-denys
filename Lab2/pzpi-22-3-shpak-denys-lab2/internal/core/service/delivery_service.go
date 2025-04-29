package service // import "wayra/internal/core/service"

import (
	"wayra/internal/core/domain/models"
	"wayra/internal/core/port"
	"wayra/internal/core/port/services"
)

// DeliveryService is a struct to manage the delivery service
type DeliveryService struct {
	*GenericService[models.Delivery] // Embedding the GenericService
}

// NewDeliveryService is a function to create a new DeliveryService
// repo: Repository of the Delivery
// return: A new DeliveryService
func NewDeliveryService(repo port.Repository[models.Delivery]) services.DeliveryService {
	return &DeliveryService{
		GenericService: NewGenericService(repo),
	}
}
