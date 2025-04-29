package models // import "wayra/internal/core/domain/models"

// ProductCategory represents a product category
type ProductCategory struct {
	// ID is the unique identifier of the product category
	// Example: 1
	ID uint `gorm:"primaryKey;column:id"`

	// Name is the name of the product category
	// Example: "Fruits"
	Name string `gorm:"size:255;not null;column:name"`

	// Description is the description of the product category
	// Example: "Fruits are a good source of vitamins and minerals"
	Description string `gorm:"type:text;column:description"`

	// MinTemperature is the minimum temperature that the product category can be stored
	// Example: 0
	MinTemperature float64 `gorm:"not null;column:min_temperature"`

	// MaxTemperature is the maximum temperature that the product category can be stored
	// Example: 10
	MaxTemperature float64 `gorm:"not null;column:max_temperature"`

	// MinHumidity is the minimum humidity that the product category can be stored
	// Example: 0
	MinHumidity float64 `gorm:"not null;column:min_humidity"`

	// MaxHumidity is the maximum humidity that the product category can be stored
	// Example: 100
	MaxHumidity float64 `gorm:"not null;column:max_humidity"`

	// IsPerishable is a flag that indicates if the product category is perishable
	// Example: true
	IsPerishable bool `gorm:"not null;column:is_perishable"`
}
