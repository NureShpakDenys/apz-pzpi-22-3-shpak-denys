package company

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"wayra/internal/adapter/httpserver/handlers"
	"wayra/internal/core/domain/dtos"
	"wayra/internal/core/domain/models"

	dtoMapper "github.com/dranikpg/dto-mapper"
	"github.com/gin-gonic/gin"
)

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
	userID, err := handlers.GetUserIDFromToken(c)
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
	userID, err := handlers.GetUserIDFromToken(c)
	if err != nil {
		userID = nil
	}
	h.LogService.LogAction(userID, actionType, description, success)
}
