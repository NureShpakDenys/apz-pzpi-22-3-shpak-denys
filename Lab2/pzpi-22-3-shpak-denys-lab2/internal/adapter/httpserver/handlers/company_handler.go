package handlers // import "wayra/internal/adapter/httpserver/handlers"

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"wayra/internal/core/domain/dtos"
	"wayra/internal/core/domain/models"
	"wayra/internal/core/port/services"

	dtoMapper "github.com/dranikpg/dto-mapper"
	"github.com/golang-jwt/jwt/v5"

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
	Role Role `json:"role" example:"user | admin | manager"`
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

	userID, err := getUserIDFromToken(c)
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
		Role:      string(RoleAdmin),
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

	userID, err := getUserIDFromToken(c)
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

// AddUserToCompany godoc
// @Summary      Add a user to a company
// @Description  Adds a user to a company if the request is made by the company creator
// @Tags         company
// @Accept       json
// @Produce      json
// @Param        company_id path int true "Company ID"
// @Param        userID body AddUserToCompanyRequest true "User ID to add"
// @Security     BearerAuth
// @Router       /company/{company_id}/add-user [post]
func (h *CompanyHandler) AddUserToCompany(c *gin.Context) {
	userID, err := getUserIDFromToken(c)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "you can not do this"})
		h.logAction(c, "add_user_to_company", "Failed to add user to company: "+err.Error(), false)

		return
	}

	user, err := h.UserService.GetByID(context.Background(), *userID)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "you can not do this"})
		h.logAction(c, "add_user_to_company", "Failed to add user to company: "+err.Error(), false)

		return
	}

	if user.Role.Name != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "you can not do this"})
		h.logAction(c, "add_user_to_company", "Failed to add user to company: "+err.Error(), false)

		return
	}

	var addUserToCompanyRequest AddUserToCompanyRequest

	companyID, err := strconv.Atoi(c.Param("company_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid company ID format"})
		h.logAction(c, "add_user_to_company", "Failed to add user to company: "+err.Error(), false)

		return
	}

	if err := c.ShouldBindJSON(&addUserToCompanyRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		h.logAction(c, "add_user_to_company", "Failed to add user to company: "+err.Error(), false)

		return
	}

	userCompany := models.UserCompany{
		UserID:    addUserToCompanyRequest.UserID,
		CompanyID: uint(companyID),
		Role:      string(addUserToCompanyRequest.Role),
	}

	if err := h.userCompanyService.Create(context.Background(), &userCompany); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		h.logAction(c, "add_user_to_company", "Failed to add user to company: "+err.Error(), false)

		return
	}

	userCompanyDTO := &dtos.UserCompanyDTO{}

	if err = dtoMapper.Map(userCompanyDTO, userCompany); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		h.logAction(c, "add_user_to_company", "Failed to add user to company: "+err.Error(), false)

		return
	}

	h.logAction(c, "add_user_to_company", fmt.Sprintf("User %d added to company %d", addUserToCompanyRequest.UserID, companyID), true)
	c.JSON(http.StatusOK, userCompanyDTO)
}

// UpdateUserInCompany godoc
// @Summary      Update a user in a company
// @Description  Updates a user in a company if the request is made by the company creator
// @Tags         company
// @Accept       json
// @Produce      json
// @Param        company_id path int true "Company ID"
// @Param        userID body UpdateUserInCompanyRequest true "User ID to update"
// @Security     BearerAuth
// @Router       /company/{company_id}/update-user [put]
func (h *CompanyHandler) UpdateUserInCompany(c *gin.Context) {
	var updateUserInCompanyRequest UpdateUserInCompanyRequest

	companyID, err := strconv.Atoi(c.Param("company_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid company ID format"})
		h.logAction(c, "update_user_in_company", "Failed to update user in company: "+err.Error(), false)

		return
	}

	if err := c.ShouldBindJSON(&updateUserInCompanyRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		h.logAction(c, "update_user_in_company", "Failed to update user in company: "+err.Error(), false)

		return
	}

	userCompanies, err := h.userCompanyService.Where(context.Background(), &models.UserCompany{
		UserID:    updateUserInCompanyRequest.UserID,
		CompanyID: uint(companyID),
	})
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found in the company"})
		h.logAction(c, "update_user_in_company", "Failed to update user in company: "+err.Error(), false)

		return
	}

	userCompany := &userCompanies[0]
	if updateUserInCompanyRequest.Role != "" {
		userCompany.Role = updateUserInCompanyRequest.Role
	}

	if err := h.userCompanyService.Update(context.Background(), userCompany); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		h.logAction(c, "update_user_in_company", "Failed to update user in company: "+err.Error(), false)

		return
	}

	userCompanyDTO := &dtos.UserCompanyDTO{}

	if err = dtoMapper.Map(userCompanyDTO, userCompany); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		h.logAction(c, "update_user_in_company", "Failed to update user in company: "+err.Error(), false)

		return
	}

	h.logAction(c, "update_user_in_company", fmt.Sprintf("User %d updated in company %d", updateUserInCompanyRequest.UserID, companyID), true)
	c.JSON(http.StatusOK, userCompanyDTO)
}

// RemoveUserFromCompany godoc
// @Summary      Remove a user from a company
// @Description  Removes a user from a company if the request is made by the company creator
// @Tags         company
// @Accept       json
// @Produce      json
// @Param        company_id path int true "Company ID"
// @Param        userID body RemoveUserFromCompanyRequest true "User ID to remove"
// @Security     BearerAuth
// @Router       /company/{company_id}/remove-user [delete]
func (h *CompanyHandler) RemoveUserFromCompany(c *gin.Context) {
	var removeUserFromCompanyRequest RemoveUserFromCompanyRequest

	companyID, err := strconv.Atoi(c.Param("company_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid company ID format"})
		h.logAction(c, "remove_user_from_company", "Failed to remove user from company: "+err.Error(), false)

		return
	}

	if err := c.ShouldBindJSON(&removeUserFromCompanyRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		h.logAction(c, "remove_user_from_company", "Failed to remove user from company: "+err.Error(), false)

		return
	}

	userCompanies, err := h.userCompanyService.Where(context.Background(), &models.UserCompany{
		UserID:    removeUserFromCompanyRequest.UserID,
		CompanyID: uint(companyID),
	})
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found in the company"})
		h.logAction(c, "remove_user_from_company", "Failed to remove user from company: "+err.Error(), false)

		return
	}

	userCompany := userCompanies[0]

	if err := h.userCompanyService.Delete(context.Background(), userCompany.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.logAction(c, "remove_user_from_company", fmt.Sprintf("User %d removed from company %d", removeUserFromCompanyRequest.UserID, companyID), true)
	c.JSON(http.StatusOK, gin.H{"message": "User removed from company successfully"})
}

// getUserIdFromToken gets the user ID from the token in the request
// c: The gin context
// Returns: The user ID and an error if there was a problem
// getUserIDFromToken gets the user ID from the token in the request
// c: The gin context
// Returns: The user ID and an error if there was a problem
func getUserIDFromToken(c *gin.Context) (*uint, error) {
	tokenCookie, exists := c.Get("token")
	if !exists {
		return nil, errors.New("token not found in context")
	}

	tokenString, ok := tokenCookie.(string)
	if !ok {
		return nil, errors.New("invalid token format in context")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}

		return []byte("mysecret123"), nil
	})
	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid token: %s", err.Error())
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		stringUserID, ok := claims["sub"].(string)
		if !ok {
			return nil, errors.New("sub not found in token claims")
		}

		userID, err := strconv.Atoi(stringUserID)
		if err != nil {
			return nil, errors.New("problem parsing user ID")
		}

		uintUserID := uint(userID)
		return &uintUserID, nil
	}

	return nil, errors.New("failed to parse token claims")
}

// GetCompanies godoc
// @Summary      Get all companies
// @Description  Retrieves a list of all companies
// @Tags         company
// @Produce      json
// @Security     BearerAuth
// @Router       /company/ [get]
// @Success      200  {array}  dtos.CompanyDTO
func (h *CompanyHandler) GetCompanies(c *gin.Context) {
	var companies []models.Company

	companies, err := h.companyService.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	companiesDTO := make([]dtos.CompanyDTO, len(companies))
	for i, company := range companies {
		if err = dtoMapper.Map(&companiesDTO[i], company); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}
	c.JSON(http.StatusOK, companiesDTO)
}

func (h *CompanyHandler) logAction(c *gin.Context, actionType, description string, success bool) {
	userID, err := getUserIDFromToken(c)
	if err != nil {
		userID = nil
	}
	h.LogService.LogAction(userID, actionType, description, success)
}
