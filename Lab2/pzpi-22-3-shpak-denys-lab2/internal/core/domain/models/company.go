// Package models implements the models for the domain of the application.
package models // import "wayra/internal/core/domain/models"

import "gorm.io/gorm"

// Company represents the Company entity.
// A Company is a group of users that can create routes and deliveries.
// A Company has a Creator, which is the user that created the company.
// A Company has a list of Users, which are the users that belong to the company.
// A Company has a list of Routes, which are the routes created by the company.
// A Company has a list of Deliveries, which are the deliveries created by the company.
type Company struct {
	// ID is the unique identifier of the company.
	// Example: 1
	ID uint `gorm:"primaryKey;column:id;autoIncrement"`

	// Name is the name of the company.
	// Example: "Wayra"
	Name string `gorm:"size:255;not null;column:name"`

	// Address is the address of the company.
	// Example: "Calle 123"
	Address string `gorm:"type:text;column:address"`

	// CreatorID is the ID of the user that created the company.
	// Example: 1
	CreatorID uint `gorm:"not null;column:creator_id"`

	// Creator is the user that created the company.
	Creator User `gorm:"foreignKey:CreatorID;constraint:OnDelete:CASCADE;" json:"creator,omitempty"`

	// Users is the list of users that belong to the company.
	Users []User `gorm:"many2many:user_companies;constraint:OnDelete:CASCADE;" json:"users,omitempty"`

	// Routes is the list of routes created by the company.
	Routes []Route `gorm:"foreignKey:CompanyID;constraint:OnDelete:CASCADE;" json:"routes,omitempty"`

	// Deliveries is the list of deliveries created by the company.
	Deliveries []Delivery `gorm:"foreignKey:CompanyID;constraint:OnDelete:CASCADE;" json:"deliveries,omitempty"`
}

// LoadRelations is the implementation of the LoadRelations interface for the Company model.
func (c *Company) LoadRelations(db *gorm.DB) *gorm.DB {
	withCreator := db.Preload("Creator")
	withUsers := withCreator.Preload("Users")
	withRoutes := withUsers.Preload("Routes").Preload("Routes.Waypoints").Preload("Routes.Waypoints.SensorData")
	withDeliveries := withRoutes.Preload("Deliveries").Preload("Deliveries.Products").Preload("Deliveries.Products.ProductCategory")
	return withDeliveries
}
