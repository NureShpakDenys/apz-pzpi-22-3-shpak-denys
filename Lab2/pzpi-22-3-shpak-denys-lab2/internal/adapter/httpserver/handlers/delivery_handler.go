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

// DeliveryHandler is a struct to handle delivery requests
type DeliveryHandler struct {
	deliveryService    services.DeliveryService    // Delivery service
	companyService     services.CompanyService     // Company service
	userCompanyService services.UserCompanyService // UserCompany service
	LogService         services.LogService         // Log service
}

// DeliveryStatus is a type to represent the status of a delivery
// deliveryService: Service for delivery operations
// companyService: Service for company operations
// userCompanyService: Service for user company operations
// returns: DeliveryHandler struct
func NewDeliveryHandler(
	deliveryService services.DeliveryService,
	companyService services.CompanyService,
	userCompanyService services.UserCompanyService,
	logService services.LogService,
) *DeliveryHandler {
	return &DeliveryHandler{
		deliveryService:    deliveryService,
		companyService:     companyService,
		userCompanyService: userCompanyService,
		LogService:         logService,
	}
}

// CreateDeliveryRequest is a struct to represent the request to create a delivery
type CreateDeliveryRequest struct {
	// CompanyID is the ID of the company
	// example: 1
	CompanyID uint `json:"company_id"`

	// Date is the date of the delivery
	// example: 2023-09-01
	Date string `json:"date" example:"2023-09-01"`
}

// UpdateDeliveryRequest is a struct to represent the request to update a delivery
type UpdateDeliveryRequest struct {
	// Date is the date of the delivery
	// example: 2024-08-01
	Date string `json:"date" example:"2024-08-01"`

	// Status is the status of the delivery
	// example: completed
	Status string `json:"status" example:"completed"`

	// The duration of the delivery
	// example: 2 hours 30 minutes
	Duration string `json:"duration" example:"2 hours 30 minutes"`
}

// CreateDelivery godoc
// @Summary      Create a delivery
// @Description  Create a delivery
// @Tags         delivery
// @Accept       json
// @Produce      json
// @Param        request body CreateDeliveryRequest true "CreateDeliveryRequest"
// @Security     BearerAuth
// @Router       /delivery [post]
func (h *DeliveryHandler) CreateDelivery(c *gin.Context) {
	var deliveryRequest CreateDeliveryRequest
	var err error
	if err := c.ShouldBindJSON(&deliveryRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		h.logAction(c, "create_delivery", "Failed to create delivery: "+err.Error(), false)

		return
	}

	userID, err := getUserIDFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		h.logAction(c, "create_delivery", "Failed to create delivery: "+err.Error(), false)

		return
	}

	var company *models.Company
	if company, err = h.companyService.GetByID(context.Background(), deliveryRequest.CompanyID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad company ID"})
		h.logAction(c, "create_delivery", "Failed to create delivery: "+err.Error(), false)

		return
	}

	userCompany, err := h.userCompanyService.Where(context.Background(), &models.UserCompany{
		UserID:    *userID,
		CompanyID: deliveryRequest.CompanyID,
	})
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		h.logAction(c, "create_delivery", "Failed to create delivery: "+err.Error(), false)

		return
	}

	if userCompany[0].Role != string(RoleAdmin) && userCompany[0].Role != string(RoleManager) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		h.logAction(c, "create_delivery", "Failed to create delivery: "+err.Error(), false)

		return
	}

	var date time.Time
	if deliveryRequest.Date != "" {
		parsedDate, err := time.Parse("2006-01-02", deliveryRequest.Date)
		if err != nil {
			fmt.Printf("Error parsing date: %v\n", err)
			date = time.Now()
		} else {
			date = parsedDate
		}
	}

	if len(company.Routes) <= 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Can not create dilivery because there is no routes for it"})
		return
	}

	delivery := &models.Delivery{
		Status:    NotStarted,
		Date:      date,
		Duration:  "0 hour",
		CompanyID: deliveryRequest.CompanyID,
		RouteID:   company.Routes[0].ID,
	}

	if err := h.deliveryService.Create(context.Background(), delivery); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		h.logAction(c, "create_delivery", "Failed to create delivery: "+err.Error(), false)

		return
	}

	deliveryDTO := &dtos.DeliveryDTO{}
	if err := dtoMapper.Map(deliveryDTO, delivery); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		h.logAction(c, "create_delivery", "Failed to create delivery: "+err.Error(), false)

		return
	}

	h.logAction(c, "create_delivery", fmt.Sprintf("Delivery for company %d created", deliveryRequest.CompanyID), true)
	c.JSON(http.StatusOK, deliveryDTO)
}

// GetDelivery godoc
// @Summary      Get a delivery
// @Description  Get a delivery
// @Tags         delivery
// @Accept       json
// @Produce      json
// @Param        delivery_id path int true "Delivery ID"
// @Security     BearerAuth
// @Router       /delivery/{delivery_id} [get]
func (h *DeliveryHandler) GetDelivery(c *gin.Context) {
	deliveryID, err := strconv.Atoi(c.Param("delivery_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Delivery ID format"})
		h.logAction(c, "get_delivery", "Failed to fetch delivery: "+err.Error(), false)

		return
	}

	delivery, err := h.deliveryService.GetByID(context.Background(), uint(deliveryID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		h.logAction(c, "get_delivery", "Failed to fetch delivery: "+err.Error(), false)

		return
	}

	userID, err := getUserIDFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		h.logAction(c, "get_delivery", "Failed to fetch delivery: "+err.Error(), false)

		return
	}
	if !h.userCompanyService.UserBelongsToCompany(*userID, delivery.CompanyID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}

	deliveryDTO := &dtos.DeliveryDTO{}
	if err = dtoMapper.Map(deliveryDTO, delivery); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		h.logAction(c, "get_delivery", "Failed to fetch delivery: "+err.Error(), false)

		return
	}

	h.logAction(c, "get_delivery", fmt.Sprintf("Delivery %d fetched", deliveryID), true)
	c.JSON(http.StatusOK, deliveryDTO)
}

// UpdateDelivery godoc
// @Summary      Update a delivery
// @Description  Update a delivery
// @Tags         delivery
// @Accept       json
// @Produce      json
// @Param        delivery_id path int true "Delivery ID"
// @Param        request body UpdateDeliveryRequest true "UpdateDeliveryRequest"
// @Security     BearerAuth
// @Router       /delivery/{delivery_id} [put]
func (h *DeliveryHandler) UpdateDelivery(c *gin.Context) {
	deliveryID, err := strconv.Atoi(c.Param("delivery_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Delivery ID format"})
		h.logAction(c, "update_delivery", "Failed to update delivery: "+err.Error(), false)

		return
	}

	var deliveryRequest UpdateDeliveryRequest
	if err := c.ShouldBindJSON(&deliveryRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		h.logAction(c, "update_delivery", "Failed to update delivery: "+err.Error(), false)

		return
	}

	delivery, err := h.deliveryService.GetByID(context.Background(), uint(deliveryID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		h.logAction(c, "update_delivery", "Failed to update delivery: "+err.Error(), false)

		return
	}

	userID, err := getUserIDFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		h.logAction(c, "update_delivery", "Failed to update delivery: "+err.Error(), false)

		return
	}

	userCompany, err := h.userCompanyService.Where(context.Background(), &models.UserCompany{
		UserID:    *userID,
		CompanyID: delivery.CompanyID,
	})
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		h.logAction(c, "update_delivery", "Failed to update delivery: "+err.Error(), false)

		return
	}

	if userCompany[0].Role != string(RoleAdmin) && userCompany[0].Role != string(RoleManager) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}

	if deliveryRequest.Date != "" {
		parsedDate, err := time.Parse("2006-01-02", deliveryRequest.Date)
		if err != nil {
			fmt.Printf("Error parsing date: %v\n", err)
		} else {
			delivery.Date = parsedDate
		}
	}
	if deliveryRequest.Status != "" {
		delivery.Status = deliveryRequest.Status
	}
	if deliveryRequest.Duration != "" {
		delivery.Duration = deliveryRequest.Duration
	}

	delivery.Company = models.Company{}
	delivery.Route = models.Route{}
	delivery.Products = nil

	if err := h.deliveryService.Update(context.Background(), delivery); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		h.logAction(c, "update_delivery", "Failed to update delivery: "+err.Error(), false)

		return
	}

	deliveryDTO := &dtos.DeliveryDTO{}
	if err = dtoMapper.Map(deliveryDTO, delivery); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		h.logAction(c, "update_delivery", "Failed to update delivery: "+err.Error(), false)

		return
	}

	h.logAction(c, "update_delivery", fmt.Sprintf("Delivery %d updated", deliveryID), true)

	c.JSON(http.StatusOK, deliveryDTO)
}

// DeleteDelivery godoc
// @Summary      Delete a delivery
// @Description  Delete a delivery
// @Tags         delivery
// @Accept       json
// @Produce      json
// @Param        delivery_id path int true "Delivery ID"
// @Security     BearerAuth
// @Router       /delivery/{delivery_id} [delete]
func (h *DeliveryHandler) DeleteDelivery(c *gin.Context) {
	deliveryID, err := strconv.Atoi(c.Param("delivery_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Delivery ID format"})
		h.logAction(c, "delete_delivery", "Failed to delete delivery: "+err.Error(), false)

		return
	}

	delivery, err := h.deliveryService.GetByID(context.Background(), uint(deliveryID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		h.logAction(c, "delete_delivery", "Failed to delete delivery: "+err.Error(), false)

		return
	}

	userID, err := getUserIDFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		h.logAction(c, "delete_delivery", "Failed to delete delivery: "+err.Error(), false)

		return
	}

	userCompany, err := h.userCompanyService.Where(context.Background(), &models.UserCompany{
		UserID:    *userID,
		CompanyID: delivery.CompanyID,
	})
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		h.logAction(c, "delete_delivery", "Failed to delete delivery: "+err.Error(), false)

		return
	}

	if userCompany[0].Role != string(RoleAdmin) && userCompany[0].Role != string(RoleManager) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}

	if err := h.deliveryService.Delete(context.Background(), uint(deliveryID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		h.logAction(c, "delete_delivery", "Failed to delete delivery: "+err.Error(), false)

		return
	}

	h.logAction(c, "delete_delivery", fmt.Sprintf("Delivery %d deleted", deliveryID), true)
	c.JSON(http.StatusOK, gin.H{"message": "Delivery deleted"})
}

func (h *DeliveryHandler) logAction(c *gin.Context, actionType, description string, success bool) {
	userID, err := getUserIDFromToken(c)
	if err != nil {
		userID = nil
	}
	h.LogService.LogAction(userID, actionType, description, success)
}
