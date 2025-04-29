package dtos // import "wayra/internal/core/domain/dtos"

// WaypointDTO is a data transfer object that represents a Waypoint entity
type WaypointDTO struct {
	// ID is the unique identifier of the Waypoint
	// Example: 1
	ID uint `json:"id,omitempty"`

	// Name is the name of the Waypoint
	// Example: Waypoint 1
	Name string `json:"name"`

	// DeviceSerial is the serial number of the device that sent the Waypoint
	// Example: 123456
	DeviceSerial string `json:"device_serial"`

	// Latitude is the latitude of the Waypoint
	// Example: -12.045
	Latitude float64 `json:"latitude"`

	// Longitude is the longitude of the Waypoint
	// Example: -77.0311
	Longitude float64 `json:"longitude"`

	// SendDataFrequency is the frequency of the data sent by the device in minutes
	// Example: 5
	SendDataFrequency uint `json:"send_data_frequency"`

	// GetWeatherAlerts is a boolean that indicates if the company should get weather alerts
	// Example: true
	GetWeatherAlerts bool `json:"get_weather_alerts"`

	// Status is the status of the Waypoint
	// Example: ok
	Status string `json:"status"`

	// Details is the details of the Waypoint
	// Example: Waypoint requires diagnostics
	Details string `json:"details"`

	// Altitude is the altitude of the Waypoint
	SensorData []SensorDataDTO `json:"sensor_data,omitempty"`
}
