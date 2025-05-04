package service // import "wayra/internal/core/service"

import (
	"context"
	"strings"
	"wayra/internal/core/domain/models"
	"wayra/internal/core/port"
	"wayra/internal/core/port/services"
)

// WaypointService is a service that manages waypoints
type WaypointService struct {
	*GenericService[models.Waypoint] // Embedding the generic service
}

// NewWaypointService creates a new waypoint service
// repo: the repository to use
// returns: a new waypoint service
func NewWaypointService(repo port.Repository[models.Waypoint]) services.WaypointService {
	return &WaypointService{
		GenericService: NewGenericService(repo),
	}
}

func (ws *WaypointService) Update(ctx context.Context, entity *models.Waypoint) error {
	query := "UPDATE waypoints SET "

	params := []interface{}{}
	fields := []string{}

	if entity.Name != "" {
		fields = append(fields, "name = ?")
		params = append(params, entity.Name)
	}

	if entity.DeviceSerial != "" {
		fields = append(fields, "device_serial = ?")
		params = append(params, entity.DeviceSerial)
	}

	if entity.Latitude != 0 {
		fields = append(fields, "latitude = ?")
		params = append(params, entity.Latitude)
	}

	if entity.Longitude != 0 {
		fields = append(fields, "longitude = ?")
		params = append(params, entity.Longitude)
	}

	if entity.SendDataFrequency != 0 {
		fields = append(fields, "send_data_frequency = ?")
		params = append(params, entity.SendDataFrequency)
	}

	fields = append(fields, "get_weather_alerts = ?")
	params = append(params, entity.GetWeatherAlerts)

	if entity.Status != "" {
		fields = append(fields, "status = ?")
		params = append(params, entity.Status)
	}

	if entity.Details != "" {
		fields = append(fields, "details = ?")
		params = append(params, entity.Details)
	}

	query += strings.Join(fields, ", ") + " WHERE id = ?"
	params = append(params, entity.ID)

	result := ws.Repository.DB().WithContext(ctx).Exec(query, params...)
	return result.Error
}
