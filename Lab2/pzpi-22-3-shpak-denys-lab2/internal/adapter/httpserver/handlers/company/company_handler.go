package company // import "wayra/internal/adapter/httpserver/handlers"

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"wayra/internal/adapter/httpserver/handlers"
	"wayra/internal/core/domain/dtos"
	"wayra/internal/core/domain/models"
	"wayra/internal/core/port/services"

	dtoMapper "github.com/dranikpg/dto-mapper"

	"github.com/gin-gonic/gin"
)

// CompanyHandler is a struct that handles company-related HTTP requests
type CompanyHandler struct {
	companyService     services.CompanyService     // Company service is an interface for comapny business logic
	userCompanyService services.UserCompanyService // UserCompany service is an interface for user-company business logic
	LogService         services.LogService         // Service for logging actions
	UserService        services.UserService        // Service to get user information
}

// NewCompanyHandler creates a new CompanyHandler with the provided services
// companyService: Service for company operations
// userCompanyService: Service for user-company operations
// Returns: A new CompanyHandler
func NewCompanyHandler(
	companyService services.CompanyService,
	userCompanyService services.UserCompanyService,
	logService services.LogService,
	userService services.UserService,
) *CompanyHandler {
	return &CompanyHandler{
		companyService:     companyService,
		userCompanyService: userCompanyService,
		LogService:         logService,
		UserService:        userService,
	}
}

// CompanyRequest is the request for the RegisterCompany endpoint
type CompanyRequest struct {
	// Name is the name of the company
	// Example: Wayra
	Name string `json:"name"`

	// Address is the address of the company
	// Example: 123 Main St
	Address string `json:"address"`
}

// AddUserToCompanyRequest is the request for the AddUserToCompany endpoint
type AddUserToCompanyRequest struct {
	// UserID is the ID of the user to add
	// Example: 1
	UserID uint `json:"userID"`

	// Role is the role of the user in the company
	Role handlers.Role `json:"role" example:"user | admin | manager"`
}

// Role is the role of a user in a company
type UpdateUserInCompanyRequest struct {
	// UserID is the ID of the user to update
	// Example: 1
	UserID uint `json:"userID"`

	// Role is the role of the user in the company
	Role string `json:"role" example:"user | admin | manager"`
}

// RemoveUserFromCompanyRequest is the request for the RemoveUserFromCompany endpoint
type RemoveUserFromCompanyRequest struct {
	// UserID is the ID of the user to remove
	// Example: 1
	UserID uint `json:"userID"`
}

// RegisterCompany godoc
// @Summary      Register a new company
// @Description  Registers a new company with the provided details
// @Tags         company
// @Accept       json
// @Produce      json
// @Param        company body CompanyRequest true "Company details"
// @Security     BearerAuth
// @Router       /company [post]
func (h *CompanyHandler) RegisterCompany(c *gin.Context) {
	var companyRequest CompanyRequest
	if err := c.ShouldBindJSON(&companyRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		h.logAction(c, "create_company", "Failed to create company: "+err.Error(), false)

		return
	}

	if companyRequest.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name is required"})
		h.logAction(c, "create_company", "Failed to create company: Name is required", false)

		return
	}

	if companyRequest.Address == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Address is required"})
		h.logAction(c, "create_company", "Failed to create company: Address is required", false)

		return
	}

	userID, err := handlers.GetUserIDFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		h.logAction(c, "create_company", "Failed to create company: "+err.Error(), false)

		return
	}

	company := &models.Company{
		Name:      companyRequest.Name,
		Address:   companyRequest.Address,
		CreatorID: *userID,
	}
	if err := h.companyService.Create(context.Background(), company); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		h.logAction(c, "create_company", "Failed to create company: "+err.Error(), false)

		return
	}

	if err = h.userCompanyService.Create(context.Background(), &models.UserCompany{
		UserID:    *userID,
		CompanyID: company.ID,
		Role:      string(handlers.RoleAdmin),
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		h.logAction(c, "create_company", "Failed to create company: "+err.Error(), false)
		return
	}

	companyDTO := &dtos.CompanyDTO{}

	if err = dtoMapper.Map(companyDTO, company); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		h.logAction(c, "create_company", "Failed to create company: "+err.Error(), false)
		return
	}

	h.logAction(c, "create_company", fmt.Sprintf("Company '%s' created", companyRequest.Name), true)
	c.JSON(http.StatusOK, companyDTO)
}

// GetCompany godoc
// @Summary      Get company details
// @Description  Retrieves the details of a company by its ID
// @Tags         company
// @Produce      json
// @Param        company_id path int true "Company ID"
// @Security     BearerAuth
// @Router       /company/{company_id} [get]
func (h *CompanyHandler) GetCompany(c *gin.Context) {
	companyID, err := strconv.Atoi(c.Param("company_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	userID, err := handlers.GetUserIDFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	company, err := h.companyService.GetByID(context.Background(), uint(companyID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Company not found"})
		return
	}

	if !h.userCompanyService.UserBelongsToCompany(*userID, uint(companyID)) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "You are not authorized to get this company's routes"})
		return
	}

	companyDTO := &dtos.CompanyDTO{}

	if err = dtoMapper.Map(companyDTO, company); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, companyDTO)
}

// UpdateCompany godoc
// @Summary      Update company details
// @Description  Updates the details of an existing company
// @Tags         company
// @Accept       json
// @Produce      json
// @Param        company_id path int true "Company ID"
// @Param        company body CompanyRequest true "Updated company details"
// @Security     BearerAuth
// @Router       /company/{company_id} [put]
func (h *CompanyHandler) UpdateCompany(c *gin.Context) {
	var companyRequest CompanyRequest

	companyID, err := strconv.Atoi(c.Param("company_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		h.logAction(c, "update_company", "Failed to update company: "+err.Error(), false)

		return
	}

	company, err := h.companyService.GetByID(context.Background(), uint(companyID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Company not found"})
		h.logAction(c, "update_company", "Failed to update company: "+err.Error(), false)

		return
	}

	if err := c.ShouldBindJSON(&companyRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		h.logAction(c, "update_company", "Failed to update company: "+err.Error(), false)

		return
	}

	company.Name = companyRequest.Name
	company.Address = companyRequest.Address
	company.Users = nil
	company.Routes = nil
	company.Deliveries = nil
	company.Creator = models.User{}

	if err := h.companyService.Update(context.Background(), company); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		h.logAction(c, "update_company", "Failed to update company: "+err.Error(), false)

		return
	}

	companyDTO := &dtos.CompanyDTO{}
	if err = dtoMapper.Map(companyDTO, company); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		h.logAction(c, "update_company", "Failed to update company: "+err.Error(), false)

		return
	}

	h.logAction(c, "update_company", fmt.Sprintf("Company ID %d updated", companyID), true)
	c.JSON(http.StatusOK, companyDTO)
}

// DeleteCompany godoc
// @Summary      Delete a company
// @Description  Deletes a company by its ID
// @Tags         company
// @Produce      json
// @Param        company_id path int true "Company ID"
// @Security     BearerAuth
// @Router       /company/{company_id} [delete]
func (h *CompanyHandler) DeleteCompany(c *gin.Context) {
	companyID, err := strconv.Atoi(c.Param("company_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		h.logAction(c, "delete_company", "Failed to delete company: "+err.Error(), false)

		return
	}

	if err := h.companyService.Delete(context.Background(), uint(companyID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		h.logAction(c, "delete_company", "Failed to delete company: "+err.Error(), false)

		return
	}

	h.logAction(c, "delete_company", fmt.Sprintf("Company ID %d deleted", companyID), true)
	c.JSON(http.StatusOK, gin.H{"message": "Company deleted successfully"})
}
