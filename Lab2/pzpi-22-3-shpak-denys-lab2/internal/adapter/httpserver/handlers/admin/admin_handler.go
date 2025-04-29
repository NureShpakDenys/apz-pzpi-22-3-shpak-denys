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
	dbPassword      string                // Password for the database
	userService     services.UserService  // Service for user operations
	encryptionKey   []byte                // Key for encryption
	db              *gorm.DB              //db for health check
	serverStartTime time.Time             // Server start time for uptime calculation
	LogService      services.LogService   // Service for logging actions
	AdminService    services.AdminService // Service for admin operations
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
		dbPassword:      cfg.DBPassword,
		userService:     userService,
		encryptionKey:   []byte(cfg.EncryptionKey),
		db:              db,
		serverStartTime: time.Now(),
		LogService:      logService,
		AdminService:    adminService,
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

// HealthCheckHandler godoc
// @Summary Health check
// @Description Health check endpoint
// @Security     BearerAuth
// @Tags admin
// @Accept json
// @Produce json
// @Router /admin/health [get]
func (h *AdminHandler) HealthCheckHandler(c *gin.Context) {
	userID, err := handlers.GetUserIDFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	currentUser, err := h.userService.GetByID(context.Background(), *userID)
	if err != nil || currentUser.Role.Name != "system_admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	dbStatus := "ok"
	if err := h.db.Raw("SELECT 1").Error; err != nil {
		dbStatus = "error"
	}

	c.JSON(http.StatusOK, gin.H{
		"server_time": time.Now(),
		"db_status":   dbStatus,
		"uptime":      time.Since(h.serverStartTime).String(),
	})
}

// GetServerLogs godoc
// @Summary Get server logs
// @Description Get server logs
// @Tags admin
// @Accept json
// @Produce json
// @Security     BearerAuth
// @Router /admin/logs [get]
func (ah *AdminHandler) GetServerLogs(c *gin.Context) {
	userID, err := handlers.GetUserIDFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	currentUser, err := ah.userService.GetByID(context.Background(), *userID)
	if err != nil || currentUser.Role.Name != "system_admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	logs, err := ah.LogService.Where(context.Background(), &models.Log{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve logs"})
		return
	}

	c.JSON(http.StatusOK, logs)
}

// GetSystemConfigs godoc
// @Summary Get system configs
// @Description Get system configs
// @Tags admin
// @Accept json
// @Produce json
// @Security     BearerAuth
// @Router /admin/system-configs [get]
func (ah *AdminHandler) GetSystemConfigs(c *gin.Context) {
	userID, err := handlers.GetUserIDFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	currentUser, err := ah.userService.GetByID(context.Background(), *userID)
	if err != nil || currentUser.Role.Name != "system_admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	configs := ah.AdminService.GetSystemConfig()
	c.JSON(http.StatusOK, configs)
}

type UpdateSystemConfigsRequest struct {
	TimeOutSec int `json:"timeout_sec"`
	TokenTTL   int `json:"token_ttl"`
}

// UpdateSystemConfigs godoc
// @Summary Update system configs
// @Description Update system configs
// @Tags admin
// @Accept json
// @Produce json
// @Param UpdateSystemConfigsRequest body UpdateSystemConfigsRequest true "UpdateSystemConfigsRequest"
// @Security     BearerAuth
// @Router /admin/system-configs [put]
func (ah *AdminHandler) UpdateSystemConfigs(c *gin.Context) {
	var req UpdateSystemConfigsRequest
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
	if err != nil || currentUser.Role.Name != "system_admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	if err := ah.AdminService.UpdateSystemConfig(req.TimeOutSec, req.TokenTTL); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update system configs"})
		return
	}

	configs := ah.AdminService.GetSystemConfig()

	c.JSON(http.StatusOK, configs)
}


type ClearLogsRequest struct {
	Days int `json:"days"`
}

// ClearLogs godoc
// @Summary Clear logs
// @Description Clear logs older than a specified number of days
// @Tags admin
// @Accept json
// @Produce json
// @Param ClearLogsRequest body ClearLogsRequest true "ClearLogsRequest"
// @Security     BearerAuth
// @Router /admin/clear-logs [post]
func (ah *AdminHandler) ClearLogs(c *gin.Context) {
	var req ClearLogsRequest
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
	if err != nil || currentUser.Role.Name != "system_admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	if err := ah.AdminService.ClearOldLogs(context.Background(), req.Days); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clear logs"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Logs cleared successfully"})
}
