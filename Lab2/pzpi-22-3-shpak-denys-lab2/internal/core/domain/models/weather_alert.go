package models // import "wayra/internal/core/domain/models"

// WeatherAlert struct
type WeatherAlert struct {
	Type    string `json:"type"`    // what is the type of alert
	Message string `json:"message"` // what is the message of the alert
	Details string `json:"details"` // details about the alert
}
