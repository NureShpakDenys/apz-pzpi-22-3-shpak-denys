package service // import "wayra/internal/core/service"

import (
	"wayra/internal/core/domain/models"
	"wayra/internal/core/port"
		"wayra/internal/core/port/services"
)

// ProductService is a struct that defines the service for the product entity
type ProductService struct {
	*GenericService[models.Product] // Embedding the generic service
}

// NewProductService is a function that creates a new product service
// repo: Repository of the product entity
// returns: A new product service
func NewProductService(repo port.Repository[models.Product]) services.ProductService {
	return &ProductService{
		GenericService: NewGenericService(repo),
	}
}
