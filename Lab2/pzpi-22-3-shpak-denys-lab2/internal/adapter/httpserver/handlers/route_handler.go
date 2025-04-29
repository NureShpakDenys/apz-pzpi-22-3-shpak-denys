package handlers // import "wayra/internal/adapter/httpserver/handlers"

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"
	"wayra/internal/core/domain/dtos"
	"wayra/internal/core/domain/models"
	"wayra/internal/core/port/services"

	dtoMapper "github.com/dranikpg/dto-mapper"
	"github.com/gin-gonic/gin"
)

// RouteHandler is a handler for route related requests
type RouteHandler struct {
	routeService       services.RouteService       // service to handle route related operations
	companyService     services.CompanyService     // service to handle company related operations
	userCompanyService services.UserCompanyService // service to handle user-company related operations
	deliveryService    services.DeliveryService    // service to handle delivery related operations
	LogService         services.LogService         // service to handle log related operations
}

// NewRoutesHandler creates a new RouteHandler
// routeService: service to handle route related operations
// companyService: service to handle company related operations
// userCompanyService: service to handle user-company related operations
// deliveryService: service to handle delivery related operations
// returns: a new RouteHandler
func NewRoutesHandler(
	routeService services.RouteService,
	companyService services.CompanyService,
	userCompanyService services.UserCompanyService,
	deliveryService services.DeliveryService,
	logService services.LogService,
) *RouteHandler {
	return &RouteHandler{
		routeService:       routeService,
		companyService:     companyService,
		userCompanyService: userCompanyService,
		deliveryService:    deliveryService,
		LogService:         logService,
	}
}

// CreateRouteRequest represents the request body for creating a route
type CreateRouteRequest struct {
	// Name of the route
	// example: Route 1
	Name string `json:"name" example:"Route 1"`

	// ID of the company to which the route belongs
	// example: 1
	CompanyID uint `json:"company_id"`
}

// UpdateRouteRequest represents the request body for updating a route
type UpdateRouteRequest struct {
	// Name of the route
	// example: Route 1
	Name string `json:"name" example:"Route 1"`

	// Status of the route
	// example: Normal temperature
	Status string `json:"status" example:"Normal temperature"`

	// Details of the route
	// example: Everything is fine
	Details string `json:"details" example:"Everything is fine"`
}

// CreateRoute godoc
// @Summary      Create a new route
// @Description  Creates a new route with the provided details
// @Security     BearerAuth
// @Tags         route
// @Accept       json
// @Produce      json
// @Param        route body CreateRouteRequest true "Route details"
// @Router       /routes [post]
func (h *RouteHandler) CreateRoute(c *gin.Context) {
	var routeRequest CreateRouteRequest
	if err := c.ShouldBindJSON(&routeRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		h.logAction(c, "create_route", "Failed to create route: "+err.Error(), false)

		return
	}

	userID, err := GetUserIDFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		h.logAction(c, "create_route", "Failed to create route: "+err.Error(), false)

		return
	}

	if _, err := h.companyService.GetByID(context.Background(), uint(routeRequest.CompanyID)); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	userCompany, err := h.userCompanyService.Where(context.Background(), &models.UserCompany{
		UserID:    *userID,
		CompanyID: routeRequest.CompanyID,
	})
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		h.logAction(c, "create_route", "Failed to create route: "+err.Error(), false)

		return
	}

	if userCompany[0].Role != string(RoleAdmin) && userCompany[0].Role != string(RoleManager) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}

	route := &models.Route{
		Name:      routeRequest.Name,
		CompanyID: uint(routeRequest.CompanyID),
		Status:    "Normal temperature",
		Details:   "Everything is fine",
		Company:   models.Company{},
		Waypoints: []models.Waypoint{},
	}

	if err := h.routeService.Create(context.Background(), route); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		h.logAction(c, "create_route", "Failed to create route: "+err.Error(), false)

		return
	}

	routeDTO := &dtos.RouteDTO{}
	if err := dtoMapper.Map(routeDTO, route); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		h.logAction(c, "create_route", "Failed to create route: "+err.Error(), false)

		return
	}

	h.logAction(c, "create_route", fmt.Sprintf("Route '%s' created successfully", routeRequest.Name), true)
	c.JSON(http.StatusOK, routeDTO)
}

// GetRoute godoc
// @Summary      Get a route
// @Description  Retrieves a route with the given ID
// @Tags         route
// @Produce      json
// @Param        route_id path int true "Route ID"
// @Security     BearerAuth
// @Router       /routes/{route_id} [get]
func (h *RouteHandler) GetRoute(c *gin.Context) {
	routeID, err := strconv.Atoi(c.Param("route_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Route ID format"})
		h.logAction(c, "get_route", "Failed to fetch route: "+err.Error(), false)

		return
	}

	route, err := h.routeService.GetByID(context.Background(), uint(routeID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Route not found"})
		h.logAction(c, "get_route", "Failed to fetch route: "+err.Error(), false)

		return
	}

	userID, err := GetUserIDFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		h.logAction(c, "get_route", "Failed to fetch route: "+err.Error(), false)

		return
	}

	if !h.userCompanyService.UserBelongsToCompany(*userID, route.CompanyID) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "You are not authorized to get this company's routes"})
		h.logAction(c, "get_route", "Failed to fetch route: "+err.Error(), false)

		return
	}

	routeDTO := &dtos.RouteDTO{}
	if err = dtoMapper.Map(routeDTO, route); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		h.logAction(c, "get_route", "Failed to fetch route: "+err.Error(), false)

		return
	}
	h.logAction(c, "get_route", fmt.Sprintf("Route %d fetched successfully", routeID), true)

	c.JSON(http.StatusOK, routeDTO)
}

// UpdateRoute godoc
// @Summary      Update an existing route
// @Description  Updates an existing route with the given ID
// @Tags         route
// @Accept       json
// @Produce      json
// @Param        route_id path int true "Route ID"
// @Param        route body UpdateRouteRequest true "Updated route details"
// @Security     BearerAuth
// @Router       /routes/{route_id} [put]
func (h *RouteHandler) UpdateRoute(c *gin.Context) {
	routeID, err := strconv.Atoi(c.Param("route_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Route ID format"})
		h.logAction(c, "update_route", "Failed to update route: "+err.Error(), false)

		return
	}

	route, err := h.routeService.GetByID(context.Background(), uint(routeID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Route not found or does not belong to the specified company"})
		h.logAction(c, "update_route", "Failed to update route: "+err.Error(), false)

		return
	}

	userID, err := GetUserIDFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		h.logAction(c, "update_route", "Failed to update route: "+err.Error(), false)

		return
	}

	userCompany, err := h.userCompanyService.Where(context.Background(), &models.UserCompany{
		UserID:    *userID,
		CompanyID: route.CompanyID,
	})
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		h.logAction(c, "update_route", "Failed to update route: "+err.Error(), false)

		return
	}

	if userCompany[0].Role != string(RoleAdmin) && userCompany[0].Role != string(RoleManager) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		h.logAction(c, "update_route", "Failed to update route: "+err.Error(), false)

		return
	}

	var routeRequest UpdateRouteRequest
	if err := c.ShouldBindJSON(&routeRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		h.logAction(c, "update_route", "Failed to update route: "+err.Error(), false)

		return
	}

	if routeRequest.Name != "" {
		route.Name = routeRequest.Name
	}

	if routeRequest.Status != "" {
		route.Status = routeRequest.Status
	}

	if routeRequest.Details != "" {
		route.Details = routeRequest.Details
	}

	route.Company = models.Company{}
	route.Waypoints = nil

	if err := h.routeService.Update(context.Background(), route); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		h.logAction(c, "update_route", "Failed to update route: "+err.Error(), false)

		return
	}

	routeDTO := &dtos.RouteDTO{}
	if err = dtoMapper.Map(routeDTO, route); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		h.logAction(c, "update_route", "Failed to update route: "+err.Error(), false)

		return
	}

	h.logAction(c, "update_route", fmt.Sprintf("Route %d updated successfully", routeID), true)

	c.JSON(http.StatusOK, routeDTO)
}

// DeleteRoute godoc
// @Summary      Delete a route
// @Description  Deletes a route with the given ID
// @Tags         route
// @Produce      json
// @Param        route_id path int true "Route ID"
// @Security     BearerAuth
// @Router       /routes/{route_id} [delete]
func (h *RouteHandler) DeleteRoute(c *gin.Context) {
	routeID, err := strconv.Atoi(c.Param("route_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid route ID format"})
		h.logAction(c, "delete_route", "Failed to delete route: "+err.Error(), false)

		return
	}

	route, err := h.routeService.GetByID(context.Background(), uint(routeID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Route not found or does not belong to the specified company"})
		h.logAction(c, "delete_route", "Failed to delete route: "+err.Error(), false)

		return
	}

	userID, err := GetUserIDFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		h.logAction(c, "delete_route", "Failed to delete route: "+err.Error(), false)

		return
	}

	userCompany, err := h.userCompanyService.Where(context.Background(), &models.UserCompany{
		UserID:    *userID,
		CompanyID: route.CompanyID,
	})
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		h.logAction(c, "delete_route", "Failed to delete route: "+err.Error(), false)

		return
	}

	if userCompany[0].Role != string(RoleAdmin) && userCompany[0].Role != string(RoleManager) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		h.logAction(c, "delete_route", "Failed to delete route: "+err.Error(), false)

		return
	}

	if err := h.routeService.Delete(context.Background(), uint(routeID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		h.logAction(c, "delete_route", "Failed to delete route: "+err.Error(), false)

		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Route deleted successfully"})
	h.logAction(c, "delete_route", fmt.Sprintf("Route %d deleted successfully", routeID), true)
}

// GetOptimalRoute godoc
// @Summary      Get optimal route
// @Description  Retrieves the optimal route for the given route ID
// @Tags         analytics
// @Produce      json
// @Param        delivery_id path int true "delivery_id"
// @Security     BearerAuth
// @Router       /analytics/{delivery_id}/optimal-route [get]
func (h *RouteHandler) GetOptimalRoute(c *gin.Context) {
	deliveryID, err := strconv.Atoi(c.Param("delivery_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid company ID format"})
		h.logAction(c, "get_optimal_route", "Failed to get optimal route: "+err.Error(), false)

		return
	}

	delivery, err := h.deliveryService.GetByID(context.Background(), uint(deliveryID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Delivery not found"})
		h.logAction(c, "get_optimal_route", "Failed to get optimal route: "+err.Error(), false)

		return
	}

	userID, err := GetUserIDFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		h.logAction(c, "get_optimal_route", "Failed to get optimal route: "+err.Error(), false)

		return
	}

	if !h.userCompanyService.UserBelongsToCompany(*userID, delivery.CompanyID) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "You are not authorized to get this company's routes"})
		h.logAction(c, "get_optimal_route", "Failed to get optimal route: "+err.Error(), false)

		return
	}

	message, predictData, coeffs, route, err := h.routeService.GetOptimalRoute(context.Background(), delivery, true, true)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		h.logAction(c, "get_optimal_route", "Failed to get optimal route: "+err.Error(), false)

		return
	}

	type RouteDTO struct {
		ID   uint   `json:"id"`
		Name string `json:"name"`
	}

	routeDTO := &RouteDTO{}
	if err = dtoMapper.Map(routeDTO, route); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		h.logAction(c, "get_optimal_route", "Failed to get optimal route: "+err.Error(), false)

		return
	}

	equation := fmt.Sprintf(
		"y = %f + %f * Temperature + %f * Humidity + %f * WindSpeed + %f * TotalWeight",
		coeffs[0],
		coeffs[1],
		coeffs[2],
		coeffs[3],
		coeffs[4],
	)

	h.logAction(c, "get_optimal_route", fmt.Sprintf("Optimal route calculated for delivery %d", deliveryID), true)

	c.JSON(http.StatusOK, gin.H{
		"message":      message,
		"route":        routeDTO,
		"predict_data": predictData,
		"equation":     equation,
	})
}

// GetWeatherAlert godoc
// @Summary      Get weather alert
// @Description  Retrieves the weather alert for the given route ID
// @Tags         analytics
// @Produce      json
// @Param        route_id path int true "route_id"
// @Security     BearerAuth
// @Router       /routes/{route_id}/weather-alert [get]
func (h *RouteHandler) GetWeatherAlert(c *gin.Context) {
	routeID, err := strconv.Atoi(c.Param("route_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid route ID format"})
		h.logAction(c, "get_weather_alert", "Failed to get weather alert: "+err.Error(), false)

		return
	}

	route, err := h.routeService.GetByID(context.Background(), uint(routeID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Route not found"})
		h.logAction(c, "get_weather_alert", "Failed to get weather alert: "+err.Error(), false)

		return
	}

	userID, err := GetUserIDFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		h.logAction(c, "get_weather_alert", "Failed to get weather alert: "+err.Error(), false)

		return
	}

	if !h.userCompanyService.UserBelongsToCompany(*userID, uint(route.CompanyID)) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "You are not authorized to access this route"})
		h.logAction(c, "get_weather_alert", "Failed to get weather alert: "+err.Error(), false)

		return
	}

	alerts, err := h.routeService.GetWeatherAlert(context.Background(), *route)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		h.logAction(c, "get_weather_alert", "Failed to get weather alert: "+err.Error(), false)

		return
	}

	h.logAction(c, "get_weather_alert", fmt.Sprintf("Weather alert fetched for route %d", routeID), true)

	c.JSON(http.StatusOK, gin.H{"alerts": alerts})
}

// GetOptimalBackRoute godoc
// @Summary      Get optimal back route
// @Description  Retrieves the optimal back route for the given route ID
// @Tags         analytics
// @Produce      json
// @Param        delivery_id path int true "delivery_id"
// @Security     BearerAuth
// @Router       /analytics/{delivery_id}/optimal-back-route [get]
func (h *RouteHandler) GetOptimalBackRoute(c *gin.Context) {
	deliveryID, err := strconv.Atoi(c.Param("delivery_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid company ID format"})
		h.logAction(c, "get_optimal_back_route", "Failed to get optimal back route: "+err.Error(), false)

		return
	}

	delivery, err := h.deliveryService.GetByID(context.Background(), uint(deliveryID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Delivery not found"})
		h.logAction(c, "get_optimal_back_route", "Failed to get optimal back route: "+err.Error(), false)

		return
	}

	userID, err := GetUserIDFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		h.logAction(c, "get_optimal_back_route", "Failed to get optimal back route: "+err.Error(), false)

		return
	}

	if !h.userCompanyService.UserBelongsToCompany(*userID, delivery.CompanyID) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "You are not authorized to get this company's routes"})
		return
	}

	message, predictData, coeffs, route, err := h.routeService.GetOptimalRoute(context.Background(), delivery, false, false)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		h.logAction(c, "get_optimal_back_route", "Failed to get optimal back route: "+err.Error(), false)

		return
	}

	type RouteDTO struct {
		ID   uint   `json:"id"`
		Name string `json:"name"`
	}

	routeDTO := &RouteDTO{}
	if err = dtoMapper.Map(routeDTO, route); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		h.logAction(c, "get_optimal_back_route", "Failed to get optimal back route: "+err.Error(), false)

		return
	}

	equation := fmt.Sprintf(
		"y = %f + %f * Temperature + %f * Humidity + %f * WindSpeed",
		coeffs[0],
		coeffs[1],
		coeffs[2],
		coeffs[3],
	)

	h.logAction(c, "get_optimal_back_route", fmt.Sprintf("Optimal back route calculated for delivery %d", deliveryID), true)

	c.JSON(http.StatusOK, gin.H{
		"message":      message,
		"route":        routeDTO,
		"predict_data": predictData,
		"equation":     equation,
	})
}

// GetSensorData godoc
// @Summary      Get sensor data
// @Description  Retrieves the sensor data for the given route ID
// @Tags         routes
// @Produce      json
// @Param        route_id path int true "route_id"
// @Param        for_hours query int false "Number of hours to get data for"
// @Router       /routes/{route_id}/get-sensor-data [get]
func (h *RouteHandler) GetSensorData(c *gin.Context) {
	routeID, err := strconv.Atoi(c.Param("route_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid route ID format"})
		h.logAction(c, "get_sensor_data", "Failed to get sensor data: "+err.Error(), false)

		return
	}

	route, err := h.routeService.GetByID(context.Background(), uint(routeID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Route not found"})
		h.logAction(c, "get_sensor_data", "Failed to get sensor data: "+err.Error(), false)

		return
	}

	forHours := c.Query("for_hours")
	if forHours == "" {
		forHours = "24"
	}

	intForHours, err := strconv.Atoi(forHours)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid for_hours format"})
		h.logAction(c, "get_sensor_data", "Failed to get sensor data: "+err.Error(), false)

		return
	}

	hours := time.Now().Add(-time.Hour * time.Duration(intForHours))

	var sensorDatas []models.SensorData
	for _, waypoint := range route.Waypoints {
		waypointSensorData := waypoint.SensorData

		for _, sensorData := range waypointSensorData {
			if sensorData.Date.After(hours) {
				sensorDatas = append(sensorDatas, sensorData)
			}
		}

	}

	var sensorDataDTOs []dtos.SensorDataDTO
	if err = dtoMapper.Map(&sensorDataDTOs, sensorDatas); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		h.logAction(c, "get_sensor_data", "Failed to get sensor data: "+err.Error(), false)

		return
	}

	h.logAction(c, "get_sensor_data", fmt.Sprintf("Sensor data fetched for route %d", routeID), true)

	c.JSON(http.StatusOK, gin.H{"sensor_data": sensorDatas})
}

func (h *RouteHandler) logAction(c *gin.Context, actionType, description string, success bool) {
	userID, err := GetUserIDFromToken(c)
	if err != nil {
		userID = nil
	}
	h.LogService.LogAction(userID, actionType, description, success)
}
