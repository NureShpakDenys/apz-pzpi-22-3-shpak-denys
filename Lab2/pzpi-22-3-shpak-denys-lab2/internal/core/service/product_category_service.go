package service // import "wayra/internal/core/service"

import (
	"wayra/internal/core/domain/models"
	"wayra/internal/core/port"
	"wayra/internal/core/port/services"
)

// ProductCategoryService is a struct to manage the product category service
type ProductCategoryService struct {
	*GenericService[models.ProductCategory] // Embedding the generic service
}

// NewProductCategoryService is a function to create a new product category service
// repo: port.Repository[models.ProductCategory] is the repository to use
// returns: *ProductCategoryService
func NewProductCategoryService(repo port.Repository[models.ProductCategory]) services.ProductCategoryService {
	return &ProductCategoryService{
		GenericService: NewGenericService(repo),
	}
}
