package services // import "wayra/internal/core/port/services"

import (
	"wayra/internal/core/domain/models"
)

// ProductCategoryService is an interface to represent the product category service
type ProductCategoryService interface {
	Service[models.ProductCategory]
}
