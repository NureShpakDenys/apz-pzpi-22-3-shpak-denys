// Package dtos defines the data transfer objects used in the application.
package dtos // import "wayra/internal/core/domain/dtos"

// CompanyDTO represents the data transfer object for the company entity.
type CompanyDTO struct {
	// ID is the unique identifier of the company.
	// Example: 1
	ID uint `json:"id"`

	// Name is the name of the company.
	// Example: "TechTeam"
	Name string `json:"name"`

	// Address is the address of the company.
	// Example: "Koningin Wilhelminaplein 13, 1062 HH Amsterdam"
	Address string `json:"address"`

	// Creator is the user who created the company.
	Creator UserDTO `json:"creator"`

	// Users is the list of users that belong to the company.
	Users []UserDTO `json:"users"`

	// Routes is the list of routes that belong to the company.
	Routes []RouteDTO `json:"routes"`

	// Deliveries is the list of deliveries that belong to the company.
	Deliveries []DeliveryDTO `json:"deliveries"`
}
