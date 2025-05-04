package models // import "wayra/internal/core/domain/models"

import (
	"time"

	"gorm.io/gorm"
)

// Delivery represents a delivery entity
type Delivery struct {
	// ID represents the unique identifier of the delivery
	// Example: 1
	ID uint `gorm:"primaryKey;column:id"`

	// Status represents the status of the delivery
	// Example: "pending"
	Status string `gorm:"size:50;column:status"`

	// Date represents the date of the delivery
	// Example: "2021-08-01"
	Date time.Time `gorm:"type:timestamp;not null;column:date"`

	// Duration represents the duration of the delivery
	// Example: "2 hours"
	Duration string `gorm:"type:interval;not null;column:duration"`

	// CompanyID represents the foreign key of the company
	// Example: 1
	CompanyID uint `gorm:"not null;column:company_id"`

	// RouteID represents the foreign key of the route
	// Example: 1
	RouteID uint `gorm:"not null;column:route_id"`

	// Route represents the route of the delivery
	Route Route `gorm:"foreignKey:RouteID;constraint:OnDelete:CASCADE;" json:"route,omitempty"`

	// Company represents the company of the delivery
	Company Company `gorm:"foreignKey:CompanyID;constraint:OnDelete:CASCADE;" json:"company,omitempty"`

	// Products represents the products of the delivery
	Products []Product `gorm:"foreignKey:DeliveryID;constraint:OnDelete:CASCADE;" json:"products"`
}

// LoadRelations is an implementation of the LoadRelations interface
func (d *Delivery) LoadRelations(db *gorm.DB) *gorm.DB {
	withCompany := db.Preload("Company").Preload("Company.Creator").Preload("Company.Users")
	withRoute := withCompany.Preload("Route").Preload("Route.Waypoints")
	withProducts := withRoute.Preload("Products").Preload("Products.ProductCategory")
	return withProducts
}
