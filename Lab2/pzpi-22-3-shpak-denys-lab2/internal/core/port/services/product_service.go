package services // import "wayra/internal/core/port/services"

import (
	"wayra/internal/core/domain/models"
)

// ProductService is the interface that wraps the basic methods for the product service
type ProductService interface {
	Service[models.Product]
}
