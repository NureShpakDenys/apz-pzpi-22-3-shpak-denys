package services // import "wayra/internal/core/port/services"

import (
	"wayra/internal/core/domain/models"
)

// DeliveryService is a service that manages the delivery domain model
type DeliveryService interface {
	Service[models.Delivery]
}
