package models // import "wayra/internal/core/domain/models"

import "gorm.io/gorm"

// Product is a struct that represent the product model
type Product struct {
	// ID is the product identifier
	// Example: 1
	ID uint `gorm:"primaryKey;column:id"`

	// Name is the product name
	// Example: "Product 1"
	Name string `gorm:"size:255;not null;column:name"`

	// Weight is the product weight
	// Example: 1.5
	Weight float64 `gorm:"not null;column:weight"`

	// ProductCategoryID is the product category identifier
	// Example: 1
	ProductCategoryID uint `gorm:"not null;column:product_category_id"`

	// DeliveryID is the delivery identifier
	// Example: 1
	DeliveryID uint `gorm:"not null;column:delivery_id"`

	// ProductCategory is the product category model
	ProductCategory ProductCategory `gorm:"foreignKey:ProductCategoryID" json:"product_category,omitempty"`

	// Delivery is the delivery model
	Delivery Delivery `gorm:"foreignKey:DeliveryID;constraint:OnDelete:CASCADE;" json:"delivery,omitempty"`
}

// LoadRelations is an implementation of the interface method for the gorm.DB
func (d *Product) LoadRelations(db *gorm.DB) *gorm.DB {
	return db.Preload("ProductCategory").Preload("Delivery")
}
