package handlers // import "wayra/internal/adapter/httpserver/handlers"
import (
	"context"
	"net/http"
	"strconv"
	"wayra/internal/core/domain/models"
	"wayra/internal/core/port/services"

	"github.com/gin-gonic/gin"
)

type UserCompanyHandler struct {
	userCompanyService services.UserCompanyService // UserCompany service is an interface for user-company business logic
}

func NewUserCompanyHandler(
	userCompanyService services.UserCompanyService,
) *UserCompanyHandler {
	return &UserCompanyHandler{
		userCompanyService: userCompanyService,
	}
}

// GetCompanyUsers godoc
// @Summary      Get all users of a company
// @Description  Get all users of a company by company ID
// @Security     BearerAuth
// @Tags         company
// @Accept       json
// @Produce      json
// @Param        company_id path int true "Company ID"
// @Router       /company/{company_id}/users [get]
func (h *UserCompanyHandler) GetCompanyUsers(c *gin.Context) {
	companyID, err := strconv.Atoi(c.Param("company_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userCompanies, err := h.userCompanyService.Where(context.Background(), &models.UserCompany{CompanyID: uint(companyID)})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, userCompanies)
}
