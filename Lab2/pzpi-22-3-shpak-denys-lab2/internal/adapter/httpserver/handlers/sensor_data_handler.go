package handlers // import "wayra/internal/adapter/httpserver/handlers"

import (
	"context"
	"net/http"
	"strconv"
	"time"
	"wayra/internal/core/domain/dtos"
	"wayra/internal/core/domain/models"
	"wayra/internal/core/port/services"

	dtoMapper "github.com/dranikpg/dto-mapper"
	"github.com/gin-gonic/gin"
)

// SensorDataHandler represents the handler for managing SensorData
type SensorDataHandler struct {
	sensorDataService  services.SensorDataService  // is the service for managing SensorData
	waypointService    services.WaypointService    // is the service for managing Waypoints
	userCompanyService services.UserCompanyService // is the service for managing UserCompany
}

// NewSensorDataHandler creates a new SensorDataHandler
// sensorDataService: is the service for managing SensorData
// waypointService: is the service for managing Waypoints
// userCompanyService: is the service for managing UserCompany
// returns a new SensorDataHandler
func NewSensorDataHandler(
	sensorDataService services.SensorDataService,
	waypointService services.WaypointService,
	userCompanyService services.UserCompanyService,
) *SensorDataHandler {
	return &SensorDataHandler{
		sensorDataService:  sensorDataService,
		waypointService:    waypointService,
		userCompanyService: userCompanyService,
	}
}

// CreateSensorDataRequest represents the request for creating a SensorData
type CreateSensorDataRequest struct {
	// Date is the date and time when the sensor data was recorded
	// example: 2021-09-01T12:00:00Z
	Date string `json:"date" example:"2021-09-01T12:00:00Z"`

	// Temperature is the temperature recorded by the sensor
	// example: 25.5
	Temperature float64 `json:"temperature" example:"25.5"`

	// Humidity is the humidity recorded by the sensor
	// example: 50.0
	Humidity float64 `json:"humidity" example:"50.0"`

	// WindSpeed is the wind speed recorded by the sensor
	// example: 10.5
	WindSpeed float64 `json:"wind_speed" example:"10.5"`

	// MeanPressure is the mean pressure recorded by the sensor
	// example: 1013.25
	MeanPressure float64 `json:"pressure" example:"1013.25"`

	// WaypointID is the ID of the waypoint where the sensor data was recorded
	// example: 1
	WaypointID uint `json:"waypoint_id"`
}

// UpdateSensorDataRequest represents the request for updating a SensorData
type UpdateSensorDataRequest struct {
	// Date is the date and time when the sensor data was recorded
	// example: 2021-09-01T12:00:00Z
	Date string `json:"date" example:"2021-09-01T12:00:00Z"`

	// Temperature is the temperature recorded by the sensor
	// example: 25.5
	Temperature float64 `json:"temperature" example:"25.5"`

	// Humidity is the humidity recorded by the sensor
	// example: 50.0
	Humidity float64 `json:"humidity" example:"50.0"`
}

// AddSensorData godoc
// @Summary      Add sensor data to a SensorData
// @Description  Adds new sensor data to the specified SensorData
// @Tags         sensor
// @Accept       json
// @Produce      json
// @Param        sensor_data body CreateSensorDataRequest true "Sensor data details"
// @Router       /sensor-data/ [post]
func (h *SensorDataHandler) AddSensorData(c *gin.Context) {
	var sensorDataRequest CreateSensorDataRequest
	if err := c.ShouldBindJSON(&sensorDataRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	var date time.Time
	if sensorDataRequest.Date != "" {
		d, err := time.Parse(time.RFC3339, sensorDataRequest.Date)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format"})
			return
		}

		date = d
	}

	if sensorDataRequest.WaypointID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Waypoint ID is required"})
		return
	}

	if sensorDataRequest.Temperature == 0 || sensorDataRequest.Humidity == 0 ||
		sensorDataRequest.WindSpeed == 0 || sensorDataRequest.MeanPressure == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Temperature, Humidity, Wind Speed and Mean Pressure are required"})
		return
	}

	sensorData := &models.SensorData{
		Date:         date,
		Temperature:  sensorDataRequest.Temperature,
		Humidity:     sensorDataRequest.Humidity,
		WindSpeed:    sensorDataRequest.WindSpeed,
		MeanPressure: sensorDataRequest.MeanPressure,
		WaypointID:   sensorDataRequest.WaypointID,
		Waypoint:     models.Waypoint{},
	}

	if err := h.sensorDataService.Create(context.Background(), sensorData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	sensorDataDTO := &dtos.SensorDataDTO{}
	if err := dtoMapper.Map(sensorDataDTO, sensorData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, sensorDataDTO)
}

// GetSensorData godoc
// @Summary      Get sensor data by ID
// @Description  Retrieves sensor data with the given ID
// @Tags         sensor
// @Produce      json
// @Param        sensor_data_id path int true "Sensor Data ID"
// @Router       /sensor-data/{sensor_data_id} [get]
func (h *SensorDataHandler) GetSensorData(c *gin.Context) {
	sensorDataID, err := strconv.Atoi(c.Param("sensor_data_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Sensor Data ID format"})
		return
	}

	sensorData, err := h.sensorDataService.GetByID(context.Background(), uint(sensorDataID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	userID, err := getUserIDFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if !h.userCompanyService.UserBelongsToCompany(*userID, sensorData.Waypoint.Route.CompanyID) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "You are not authorized to update this sensorData"})
		return
	}

	sensorDataDTO := &dtos.SensorDataDTO{}
	if err = dtoMapper.Map(sensorDataDTO, sensorData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, sensorDataDTO)
}

// UpdateSensorData godoc
// @Summary      Update sensor data by ID
// @Description  Updates sensor data with the given ID
// @Tags         sensor
// @Accept       json
// @Produce      json
// @Param        sensor_data_id path int true "Sensor Data ID"
// @Param        sensor_data body UpdateSensorDataRequest true "Sensor data details"
// @Router       /sensor-data/{sensor_data_id} [put]
func (h *SensorDataHandler) UpdateSensorData(c *gin.Context) {
	sensorDataID, err := strconv.Atoi(c.Param("sensor_data_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Sensor Data ID format"})
		return
	}

	sensorData, err := h.sensorDataService.GetByID(context.Background(), uint(sensorDataID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	userID, err := getUserIDFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userCompany, err := h.userCompanyService.Where(context.Background(), &models.UserCompany{
		UserID:    *userID,
		CompanyID: sensorData.Waypoint.Route.CompanyID,
	})
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	if userCompany[0].Role != string(RoleAdmin) && userCompany[0].Role != string(RoleManager) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}

	var sensorDataRequest UpdateSensorDataRequest
	if err := c.ShouldBindJSON(&sensorDataRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	var date time.Time
	if sensorDataRequest.Date != "" {
		date, err = time.Parse(time.RFC3339, sensorDataRequest.Date)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format"})
			return
		}
	}

	sensorData.Date = date
	sensorData.Temperature = sensorDataRequest.Temperature
	sensorData.Humidity = sensorDataRequest.Humidity
	sensorData.Waypoint = models.Waypoint{}

	if err := h.sensorDataService.Update(context.Background(), sensorData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	sensorDataDTO := &dtos.SensorDataDTO{}
	if err = dtoMapper.Map(sensorDataDTO, sensorData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, sensorDataDTO)
}

// DeleteSensorData godoc
// @Summary      Delete sensor data by ID
// @Description  Deletes sensor data with the given ID
// @Tags         sensor
// @Produce      json
// @Param        sensor_data_id path int true "Sensor Data ID"
// @Security     BearerAuth
// @Router       /sensor-data/{sensor_data_id} [delete]
func (h *SensorDataHandler) DeleteSensorData(c *gin.Context) {
	sensorDataID, err := strconv.Atoi(c.Param("sensor_data_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Sensor Data ID format"})
		return
	}

	sensorData, err := h.sensorDataService.GetByID(context.Background(), uint(sensorDataID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Sensor Data not found"})
		return
	}

	userID, err := getUserIDFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userCompany, err := h.userCompanyService.Where(context.Background(), &models.UserCompany{
		UserID:    *userID,
		CompanyID: sensorData.Waypoint.Route.CompanyID,
	})
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	if userCompany[0].Role != string(RoleAdmin) && userCompany[0].Role != string(RoleManager) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}

	if err := h.sensorDataService.Delete(context.Background(), uint(sensorDataID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Sensor Data deleted successfully"})
}
