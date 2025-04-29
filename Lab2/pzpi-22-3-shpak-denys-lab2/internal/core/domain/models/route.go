package models // import "wayra/internal/core/domain/models"

import "gorm.io/gorm"

// Route is a model that represents a route on the delivery way
type Route struct {
	// ID is the identifier of the route
	// Example: 1
	ID uint `gorm:"primaryKey;column:id"`

	// Name is the name of the route
	// Example: "Route 1"
	Name string `gorm:"size:255;not null;column:name"`

	// CompanyID is the identifier of the company that the route belongs to
	// Example: 1
	CompanyID uint `gorm:"not null;column:company_id"`

	// Status is the status of the route
	// Example: "normal temperature"
	Status string `gorm:"size:255;column:status"`

	// Details is the details of the route
	// Example: "Route 1 details"
	Details string `gorm:"size:255;column:details"`

	// Company is the company that the route belongs to
	Company Company `gorm:"foreignKey:CompanyID" json:"company,omitempty"`

	// Waypoints is the list of waypoints that the route has
	Waypoints []Waypoint `gorm:"foreignKey:RouteID;constraint:OnDelete:CASCADE;" json:"waypoints,omitempty"`
}

// LoadRelations is an implementation of the LoadRelations method from the model interface
func (r *Route) LoadRelations(db *gorm.DB) *gorm.DB {
	return db.Preload("Company").Preload("Company.Creator").Preload("Waypoints").Preload("Waypoints.SensorData")
}
