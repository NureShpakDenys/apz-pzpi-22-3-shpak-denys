package models // import "wayra/internal/core/domain/models"

import "gorm.io/gorm"

// Waypoint is a struct that represents the Waypoint model of the database
type Waypoint struct {
	// ID is the identifier of the waypoint
	// Example: 1
	ID uint `gorm:"primaryKey;column:id"`

	// Name is the name of the waypoint
	// Example: Waypoint 1
	Name string `gorm:"size:255;not null;column:name"`

	// DeviceSerial is the serial number of the device that sent the waypoint
	// Example: 123456789
	DeviceSerial string `gorm:"size:255;not null;column:device_serial"`

	// Latitude is the latitude of the waypoint
	// Example: -12.04318
	Latitude float64 `gorm:"not null;column:latitude"`

	// Longitude is the longitude of the waypoint
	// Example: -77.02824
	Longitude float64 `gorm:"not null;column:longitude"`

	// SendDataFrequency is the frequency of the data sent by the device in minutes
	// Example: 5
	SendDataFrequency uint `gorm:"not null;column:send_data_frequency"`

	// GetWeatherAlerts is a boolean that indicates if the company should get weather alerts
	// Example: true
	GetWeatherAlerts bool `gorm:"not null;default:false;column:get_weather_alerts"`

	// Status is the status of the waypoint
	// Example: ok
	Status string `gorm:"size:255;column:status"`

	// Details is the details of the waypoint
	// Example: Waypoint requires diagnostics
	Details string `gorm:"size:255;column:details"`

	// RouteId is the identifier of the route to which the waypoint belongs
	// Example: 1
	RouteID uint `gorm:"not null;column:route_id"`

	// Route is the route to which the waypoint belongs
	Route Route `gorm:"foreignKey:RouteID" json:"route,omitempty"`

	// SensorData is the sensor data of the waypoint
	SensorData []SensorData `gorm:"foreignKey:WaypointID;constraint:OnDelete:CASCADE;" json:"sensor_data,omitempty"`
}

// LoadRelations is an implementation of the LoadRelations interface
func (w *Waypoint) LoadRelations(db *gorm.DB) *gorm.DB {
	withRoute := db.Preload("Route").Preload("Route.Company").Preload("Route.Company.Creator")
	withSensorData := withRoute.Preload("SensorData")
	return withSensorData
}
