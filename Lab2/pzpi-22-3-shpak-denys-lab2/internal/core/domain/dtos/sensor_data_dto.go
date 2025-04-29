package dtos // import "wayra/internal/core/domain/dtos"

import "time"

// SensorDataDTO is a DTO for SensorData
type SensorDataDTO struct {
	// ID is the unique identifier of the SensorData
	// Example: 1
	ID uint `json:"id,omitempty"`

	// Date is the date of the SensorData
	// Example: 2021-01-01T00:00:00Z
	Date time.Time `json:"date"`

	// Temperature is the temperature of the SensorData
	// Example: 25.5
	Temperature float64 `json:"temperature"`

	// Humidity is the humidity of the SensorData
	// Example: 0.5
	Humidity float64 `json:"humidity"`

	// Rain is the rain of the SensorData
	// Example: 0.5
	WindSpeed float64 `json:"wind_speed"`

	// WindDirection is the wind direction of the SensorData
	// Example: 0.5
	MeanPressure float64 `json:"mean_pressure"`
}
