package dtos // import "wayra/internal/core/domain/dtos"

// ProductCategoryDTO is a struct that represents the product category data transfer object
// It is used to represent the product category data in the application
type ProductCategoryDTO struct {
	// ID is the product category identifier
	// Example: 1
	ID uint `json:"id"`

	// Name is the product category name
	// Example: Food
	Name string `json:"name"`

	// Description is the product category description
	// Example: Food products
	Description string `json:"description"`

	// MinTemperature is the minimum temperature that the product category can be stored
	// Example: 10
	MinTemperature float64 `json:"min_temperature"`

	// MaxTemperature is the maximum temperature that the product category can be stored
	// Example: 20
	MaxTemperature float64 `json:"max_temperature"`

	// MinHumidity is the minimum humidity that the product category can be stored
	// Example: 10
	MinHumidity float64 `json:"min_humidity"`

	// MaxHumidity is the maximum humidity that the product category can be stored
	// Example: 20
	MaxHumidity float64 `json:"max_humidity"`

	// IsPerishable is a flag that indicates if the product category is perishable
	// Example: true
	IsPerishable bool `json:"is_perishable"`
}
