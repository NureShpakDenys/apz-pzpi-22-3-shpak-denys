package admin // import "wayra/internal/adapter/httpserver/handlers/admin"

import (
	"context"
	"fmt"
	"net/http"
	"time"
	"wayra/internal/adapter/config"
	"wayra/internal/adapter/httpserver/handlers"
	"wayra/internal/core/domain/models"
	"wayra/internal/core/port/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AdminHandler is a handler for admin endpoints
type AdminHandler struct {
	userService     services.UserService  // Service for user operations
	db              *gorm.DB              //db for health check
	serverStartTime time.Time             // Server start time for uptime calculation
	LogService      services.LogService   // Service for logging actions
	AdminService    services.AdminService // Service for admin operations
	Cfg             *config.Config
}

// NewAdminHandler creates a new AdminHandler
// dbPassword: Password for the database
// userService: Service for user operations
// Returns: A new AdminHandler
func NewAdminHandler(
	userService services.UserService,
	cfg *config.Config,
	db *gorm.DB,
	logService services.LogService,
	adminService services.AdminService,
) *AdminHandler {
	return &AdminHandler{
		userService:     userService,
		db:              db,
		serverStartTime: time.Now(),
		LogService:      logService,
		AdminService:    adminService,
		Cfg:             cfg,
	}
}

type UpdateUserRoleRequest struct {
	UserID int64 `json:"userId"`
	RoleID uint  `json:"roleId"`
}

// UpdateUserRole godoc
// @Summary Update user role
// @Description Update user role
// @Tags admin
// @Accept json
// @Produce json
// @Param UpdateUserRoleRequest body UpdateUserRoleRequest true "UpdateUserRoleRequest"
// @Security     BearerAuth
// @Router /admin/change-role [post]
func (ah *AdminHandler) UpdateUserRole(c *gin.Context) {
	var req UpdateUserRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := handlers.GetUserIDFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	currentUser, err := ah.userService.GetByID(context.Background(), *userID)
	if err != nil || currentUser.Role.Name != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	if err := ah.userService.Update(context.Background(), &models.User{
		ID:     uint(req.UserID),
		RoleID: req.RoleID,
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		ah.LogService.LogAction(userID, "update_user_role", "Failed to change user role: "+err.Error(), false)
		return
	}

	ah.LogService.LogAction(userID, "update_user_role", fmt.Sprintf("User %d role changed to %d", req.UserID, req.RoleID), true)
	c.JSON(http.StatusOK, gin.H{"message": "User role updated"})
}
