package dtos // import "wayra/internal/core/domain/dtos"

import "time"

// DeliveryDTO is the data transfer object for the Delivery entity
type DeliveryDTO struct {
	// ID is the unique identifier of the delivery
	// Example: 1
	ID uint `json:"id"`

	// Status is the status of the delivery
	// Example: completed
	Status string `json:"status"`

	// Date is the date of the delivery
	// Example: 2021-01-01T00:00:00Z
	Date time.Time `json:"date"`

	// Duration is the duration of the delivery
	// Example: 1h
	Duration string `json:"duration"`

	// CompanyID is the unique identifier of the company
	// Example: 1
	CompanyID uint `json:"company_id"`

	// RouteID is the unique identifier of the route
	// Example: 1
	RouteID uint `json:"route_id"`

	// Products is the list of products of the delivery
	Products []ProductDTO `json:"products,omitempty"`
}
