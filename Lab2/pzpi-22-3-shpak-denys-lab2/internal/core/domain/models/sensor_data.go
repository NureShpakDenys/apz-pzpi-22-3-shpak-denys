package models // import "wayra/internal/core/domain/models"

import (
	"time"

	"gorm.io/gorm"
)

// SensorData is a struct that represents the sensor_data table in the database
type SensorData struct {
	// ID is the primary key of the table
	// Example: 1
	ID uint `gorm:"primaryKey;column:id"`

	// Date is the date when the data was recorded
	// Example: 2021-08-01 12:00:00
	Date time.Time `gorm:"type:timestamp;not null;column:date"`

	// Temperature is the temperature recorded by the sensor
	// Example: 25.5
	Temperature float64 `gorm:"not null;column:temperature"`

	// Humidity is the humidity recorded by the sensor
	// Example: 0.5
	Humidity float64 `gorm:"not null;column:humidity"`

	// WindSpeed is the wind speed recorded by the sensor
	// Example: 10.5
	WindSpeed float64 `gorm:"not null;column:wind_speed"`

	// MeanPressure is the mean pressure recorded by the sensor
	// Example: 1013.25
	MeanPressure float64 `gorm:"not null;column:mean_pressure"`

	// WaypointID is the foreign key of the waypoint table
	// Example: 1
	WaypointID uint `gorm:"not null;column:waypoint_id"`

	// Waypoint is the relation with the waypoint table
	Waypoint Waypoint `gorm:"foreignKey:WaypointID" json:"device,omitempty"`
}

// LoadRelations is an implementation of the interface for the gorm library
func (s *SensorData) LoadRelations(db *gorm.DB) *gorm.DB {
	return db
}
