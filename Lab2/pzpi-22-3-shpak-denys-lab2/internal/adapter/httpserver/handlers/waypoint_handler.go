package handlers // import "wayra/internal/adapter/httpserver/handlers"

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"wayra/internal/core/domain/dtos"
	"wayra/internal/core/domain/models"
	"wayra/internal/core/port/services"

	dtoMapper "github.com/dranikpg/dto-mapper"
	"github.com/gin-gonic/gin"
)

// WaypointHandler is
type WaypointHandler struct {
	waypointService    services.WaypointService
	routeService       services.RouteService
	companyService     services.CompanyService
	userCompanyService services.UserCompanyService
	LogService         services.LogService
}

// NewWaypointHandler is a constructor for WaypointHandler
// waypointService: service to handle waypoints
// routeService: service to handle routes
// companyService: service to handle companies
// userCompany: service to handle user-company relationships
// returns: a new WaypointHandler
func NewWaypointHandler(
	waypointService services.WaypointService,
	routeService services.RouteService,
	companyService services.CompanyService,
	userCompany services.UserCompanyService,
	logService services.LogService,
) *WaypointHandler {
	return &WaypointHandler{
		waypointService:    waypointService,
		routeService:       routeService,
		companyService:     companyService,
		userCompanyService: userCompany,
		LogService:         logService,
	}
}

// CreateWaypointRequest is a struct to handle the request to create a waypoint
type CreateWaypointRequest struct {
	// Name of the waypoint
	// Example: "Waypoint 1"
	Name string `json:"name"`

	// Latitude of the waypoint
	// Example: -12.04318
	Latitude float64 `json:"latitude"`

	// Longitude of the waypoint
	// Example: -77.02824
	Longitude float64 `json:"longitude"`

	// Device serial number
	// Example: "1234567890"
	DeviceSerial string `json:"device_serial"`

	// SendDataFrequency of the waypoint
	// Example: 10
	SendDataFrequency uint `json:"send_data_frequency"`

	// GetWeatherAlerts of the waypoint
	// Example: true
	GetWeatherAlerts bool `json:"get_weather_alerts"`

	// Route ID to which the waypoint belongs
	// Example: 1
	RouteID uint `json:"route_id"`
}

// UpdateWaypointRequest is a struct to handle the request to update a waypoint
type UpdateWaypointRequest struct {
	// Name of the waypoint
	// Example: "Waypoint 1"
	Name string `json:"name"`

	// Latitude of the waypoint
	// Example: -12.04318
	Latitude float64 `json:"latitude"`

	// Longitude of the waypoint
	// Example: -77.02824
	Longitude float64 `json:"longitude"`

	// Device serial number
	// Example: "1234567890"
	DeviceSerial string `json:"device_serial"`

	// SendDataFrequency of the waypoint
	// Example: 10
	SendDataFrequency uint `json:"send_data_frequency"`

	// GetWeatherAlerts of the waypoint
	// Example: true
	GetWeatherAlerts bool `json:"get_weather_alerts"`

	// Status of the waypoint
	// Example: "ok"
	Status string `json:"status"`

	// Details of the waypoint
	// Example: "Everything is fine"
	Details string `json:"details"`
}

// AddWaypoint godoc
// @Summary      Add a waypoint to a route
// @Description  Adds a new waypoint to the specified route
// @Tags         waypoint
// @Accept       json
// @Produce      json
// @Param        waypoint body CreateWaypointRequest true "Waypoint details"
// @Security     BearerAuth
// @Router       /waypoints [post]
func (h *WaypointHandler) AddWaypoint(c *gin.Context) {
	var waypointRequest CreateWaypointRequest
	if err := c.ShouldBindJSON(&waypointRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		h.logAction(c, "create_waypoint", "Failed to create waypoint: "+err.Error(), false)

		return
	}

	route, err := h.routeService.GetByID(context.Background(), waypointRequest.RouteID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Route not found"})
		h.logAction(c, "create_waypoint", "Failed to create waypoint: "+err.Error(), false)

		return
	}

	userID, err := GetUserIDFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		h.logAction(c, "create_waypoint", "Failed to create waypoint: "+err.Error(), false)

		return
	}

	userCompany, err := h.userCompanyService.Where(context.Background(), &models.UserCompany{
		UserID:    *userID,
		CompanyID: route.CompanyID,
	})
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		h.logAction(c, "create_waypoint", "Failed to create waypoint: "+err.Error(), false)

		return
	}

	if userCompany[0].Role != string(RoleAdmin) && userCompany[0].Role != string(RoleManager) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		h.logAction(c, "create_waypoint", "Failed to create waypoint: "+err.Error(), false)

		return
	}

	waypoint := &models.Waypoint{
		Name:              waypointRequest.Name,
		DeviceSerial:      waypointRequest.DeviceSerial,
		Latitude:          waypointRequest.Latitude,
		Longitude:         waypointRequest.Longitude,
		SendDataFrequency: 24,
		GetWeatherAlerts:  true,
		Status:            "Ok",
		Details:           "Everything is fine",
		RouteID:           waypointRequest.RouteID,
		Route:             models.Route{},
		SensorData:        []models.SensorData{},
	}

	if err := h.waypointService.Create(context.Background(), waypoint); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		h.logAction(c, "create_waypoint", "Failed to create waypoint: "+err.Error(), false)

		return
	}

	waypointDTO := &dtos.WaypointDTO{}
	if err = dtoMapper.Map(waypointDTO, waypoint); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		h.logAction(c, "create_waypoint", "Failed to create waypoint: "+err.Error(), false)

		return
	}

	h.logAction(c, "create_waypoint", fmt.Sprintf("Waypoint '%s' created", waypointRequest.Name), true)

	c.JSON(http.StatusOK, waypointDTO)
}

// GetWaypoint godoc
// @Summary      Get waypoint details
// @Description  Retrieves the details of a waypoint
// @Tags         waypoint
// @Produce      json
// @Param        waypoint_id path int true "Waypoint ID"
// @Security     BearerAuth
// @Router       /waypoints/{waypoint_id} [get]
func (h *WaypointHandler) GetWaypoint(c *gin.Context) {
	waypointID, err := strconv.Atoi(c.Param("waypoint_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid waypoint ID format"})
		h.logAction(c, "get_waypoint", "Failed to fetch waypoint: "+err.Error(), false)

		return
	}

	waypoint, err := h.waypointService.GetByID(context.Background(), uint(waypointID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Waypoint not found"})
		h.logAction(c, "get_waypoint", "Failed to fetch waypoint: "+err.Error(), false)

		return
	}

	userID, err := GetUserIDFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		h.logAction(c, "get_waypoint", "Failed to fetch waypoint: "+err.Error(), false)

		return
	}

	if !h.userCompanyService.UserBelongsToCompany(*userID, waypoint.Route.CompanyID) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "You are not authorized to get this company's routes"})
		h.logAction(c, "get_waypoint", "Failed to fetch waypoint: "+err.Error(), false)

		return
	}

	waypointDTO := &dtos.WaypointDTO{}
	if err = dtoMapper.Map(waypointDTO, waypoint); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		h.logAction(c, "get_waypoint", "Failed to fetch waypoint: "+err.Error(), false)

		return
	}

	h.logAction(c, "get_waypoint", fmt.Sprintf("Waypoint %d fetched", waypointID), true)
	c.JSON(http.StatusOK, waypointDTO)
}

// UpdateWaypoint godoc
// @Summary      Update waypoint details
// @Description  Updates the details of a waypoint
// @Tags         waypoint
// @Accept       json
// @Produce      json
// @Param        waypoint_id path int true "Waypoint ID"
// @Param        waypoint body UpdateWaypointRequest true "Waypoint details"
// @Router       /waypoints/{waypoint_id} [put]
func (h *WaypointHandler) UpdateWaypoint(c *gin.Context) {
	waypointID, err := strconv.Atoi(c.Param("waypoint_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid waypoint ID format"})
		return
	}

	waypoint, err := h.waypointService.GetByID(context.Background(), uint(waypointID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Waypoint not found"})
		return
	}

	userID, err := GetUserIDFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userCompany, err := h.userCompanyService.Where(context.Background(), &models.UserCompany{
		UserID:    *userID,
		CompanyID: waypoint.Route.CompanyID,
	})
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	if userCompany[0].Role != string(RoleAdmin) && userCompany[0].Role != string(RoleManager) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}

	var waypointRequest UpdateWaypointRequest
	if err := c.ShouldBindJSON(&waypointRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if waypointRequest.Name != "" {
		waypoint.Name = waypointRequest.Name
	}

	if waypointRequest.DeviceSerial != "" {
		waypoint.DeviceSerial = waypointRequest.DeviceSerial
	}

	if waypointRequest.Latitude != 0 {
		waypoint.Latitude = waypointRequest.Latitude
	}

	if waypointRequest.Longitude != 0 {
		waypoint.Longitude = waypointRequest.Longitude
	}

	if waypointRequest.SendDataFrequency != 0 {
		waypoint.SendDataFrequency = waypointRequest.SendDataFrequency
	}

	if waypointRequest.GetWeatherAlerts {
		waypoint.GetWeatherAlerts = waypointRequest.GetWeatherAlerts
	}

	if waypointRequest.Status != "" {
		waypoint.Status = waypointRequest.Status
	}

	if waypointRequest.Details != "" {
		waypoint.Details = waypointRequest.Details
	}

	waypoint.Route = models.Route{}
	waypoint.SensorData = nil

	if err := h.waypointService.Update(context.Background(), waypoint); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		h.logAction(c, "update_waypoint", "Failed to update waypoint: "+err.Error(), false)

		return
	}

	waypointDTO := &dtos.WaypointDTO{}
	if err = dtoMapper.Map(waypointDTO, waypoint); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		h.logAction(c, "update_waypoint", "Failed to update waypoint: "+err.Error(), false)

		return
	}

	h.logAction(c, "update_waypoint", fmt.Sprintf("Waypoint %d updated", waypointID), true)

	c.JSON(http.StatusOK, waypointDTO)
}

// DeleteWaypoint godoc
// @Summary      Delete waypoint
// @Description  Deletes a waypoint
// @Tags         waypoint
// @Produce      json
// @Param        waypoint_id path int true "Waypoint ID"
// @Security     BearerAuth
// @Router       /waypoints/{waypoint_id} [delete]
func (h *WaypointHandler) DeleteWaypoint(c *gin.Context) {
	waypointID, err := strconv.Atoi(c.Param("waypoint_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid waypoint ID format"})
		h.logAction(c, "delete_waypoint", "Failed to delete waypoint: "+err.Error(), false)

		return
	}

	waypoint, err := h.waypointService.GetByID(context.Background(), uint(waypointID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Waypoint not found"})
		h.logAction(c, "delete_waypoint", "Failed to delete waypoint: "+err.Error(), false)

		return
	}

	userID, err := GetUserIDFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		h.logAction(c, "delete_waypoint", "Failed to delete waypoint: "+err.Error(), false)

		return
	}

	userCompany, err := h.userCompanyService.Where(context.Background(), &models.UserCompany{
		UserID:    *userID,
		CompanyID: waypoint.Route.CompanyID,
	})
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		h.logAction(c, "delete_waypoint", "Failed to delete waypoint: "+err.Error(), false)

		return
	}

	if userCompany[0].Role != string(RoleAdmin) && userCompany[0].Role != string(RoleManager) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		h.logAction(c, "delete_waypoint", "Failed to delete waypoint: "+err.Error(), false)

		return
	}

	if err := h.waypointService.Delete(context.Background(), uint(waypointID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		h.logAction(c, "delete_waypoint", "Failed to delete waypoint: "+err.Error(), false)

		return
	}

	h.logAction(c, "delete_waypoint", fmt.Sprintf("Waypoint %d deleted", waypointID), true)

	c.JSON(http.StatusOK, gin.H{"message": "Waypoint deleted successfully"})
}

// GetWaypoints  godoc
// @Summary      Get waypoints
// @Description  Retrieves all waypoints
// @Tags         waypoint
// @Produce      json
// Param		 device_serial query string false "Device serial number"
// @Router       /waypoints [get]
func (h *WaypointHandler) GetWaypoints(c *gin.Context) {
	deviceSerial := c.Query("device_serial")

	waypoints, err := h.waypointService.Where(context.Background(), &models.Waypoint{
		DeviceSerial: deviceSerial,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		h.logAction(c, "get_waypoints", "Failed to fetch waypoints: "+err.Error(), false)

		return
	}
	if len(waypoints) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No waypoints found"})
		return
	}

	waypointID := waypoints[0].ID
	routeID := waypoints[0].RouteID

	h.logAction(c, "get_waypoints", "Waypoints fetched for device serial "+deviceSerial, true)

	c.JSON(http.StatusOK, map[string]int{
		"waypoint_id": int(waypointID),
		"route_id":    int(routeID),
	})
}

// GetDeviceConfig godoc
// @Summary      Get device configuration
// @Description  Retrieves the configuration of a device
// @Tags         waypoint
// @Produce      json
// @Param        waypoint_id path int true "Waypoint ID"
// @Router       /device-config/{waypoint_id} [get]
func (h *WaypointHandler) GetDeviceConfig(c *gin.Context) {
	waypointID, err := strconv.Atoi(c.Param("waypoint_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid waypoint ID format"})
		return
	}

	waypoint, err := h.waypointService.GetByID(context.Background(), uint(waypointID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Waypoint not found"})
		return
	}

	deviceConfig := map[string]interface{}{
		"send_data_frequency": waypoint.SendDataFrequency,
		"get_weather_alerts":  waypoint.GetWeatherAlerts,
	}

	c.JSON(http.StatusOK, deviceConfig)
}

func (h *WaypointHandler) logAction(c *gin.Context, actionType, description string, success bool) {
	userID, err := GetUserIDFromToken(c)
	if err != nil {
		userID = nil
	}
	h.LogService.LogAction(userID, actionType, description, success)
}
