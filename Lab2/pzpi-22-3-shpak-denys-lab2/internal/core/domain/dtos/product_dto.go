package dtos // import "wayra/internal/core/domain/dtos"

// ProductCategoryDTO is a DTO for ProductCategory
type ProductDTO struct {
	// ID is the unique identifier of the product
	// Example: 1
	ID uint `json:"id"`

	// Name is the name of the product
	// Example: Beer
	Name string `json:"name"`

	// Description is the description of the product
	// Example: A cold beer
	Weight float64 `json:"weight"`

	// ProductCategory is the category of the product
	ProductCategory ProductCategoryDTO `json:"product_category,omitempty"`
}
