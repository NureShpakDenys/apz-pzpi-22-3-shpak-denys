package admin

import (
	"context"
	"net/http"
	"wayra/internal/adapter/httpserver/handlers"

	"github.com/gin-gonic/gin"
)

// BackupDatabaseRequest is the request for the BackupDatabase endpoint
type BackupDatabaseRequest struct {
	// Backup path is the path where the backup will be stored
	BackupPath string `json:"backup_path"`
}

// BackupDatabase godoc
// @Summary Backup the database
// @Description Backup the database in CSV format per table
// @Tags admin
// @Accept json
// @Produce json
// @Param BackupDatabaseRequest body BackupDatabaseRequest true "Backup directory path"
// @Security     BearerAuth
// @Router /admin/backup [post]
func (ah *AdminHandler) BackupDatabase(c *gin.Context) {
	var req BackupDatabaseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := handlers.GetUserIDFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	user, err := ah.userService.GetByID(context.Background(), *userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if user.Role.Name != "db_admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	if err := ah.AdminService.BackupDatabase(c.Request.Context(), *userID, req.BackupPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Backup failed: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Backup created successfully"})
}

// RestoreDatabaseRequest is the request for the RestoreDatabase endpoint
type RestoreDatabaseRequest struct {
	// Backup path is the path where the backup is stored
	BackupPath string `json:"backup_path"`
}

// RestoreDatabase godoc
// @Summary Restore the database from backup
// @Description Restore the database from encrypted backup files
// @Tags admin
// @Accept json
// @Produce json
// @Param RestoreDatabaseRequest body RestoreDatabaseRequest true "Backup directory path"
// @Security     BearerAuth
// @Success 200 {string} string "Database restored"
// @Router /admin/restore [post]
func (ah *AdminHandler) RestoreDatabase(c *gin.Context) {
	var req RestoreDatabaseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := handlers.GetUserIDFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	user, err := ah.userService.GetByID(context.Background(), *userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if user.Role.Name != "db_admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	if err := ah.AdminService.RestoreDatabase(c.Request.Context(), *userID, req.BackupPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Restore failed: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Database restored successfully"})
}

// GetDBStatus godoc
// @Summary Get database status
// @Description Get database status
// @Tags admin
// @Accept json
// @Produce json
// @Security     BearerAuth
// @Router /admin/db-status [get]
func (ah *AdminHandler) GetDBStatus(c *gin.Context) {
	userID, err := handlers.GetUserIDFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	currentUser, err := ah.userService.GetByID(context.Background(), *userID)
	if err != nil || currentUser.Role.Name != "db_admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	dbStatus, _ := ah.AdminService.GetDBStatus(context.Background())

	c.JSON(http.StatusOK, dbStatus)
}

// OptimizeDatabase godoc
// @Summary Optimize database
// @Description Optimize database
// @Tags admin
// @Accept json
// @Produce json
// @Security     BearerAuth
// @Router /admin/optimize [post]
func (ah *AdminHandler) OptimizeDatabase(c *gin.Context) {
	userID, err := handlers.GetUserIDFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	currentUser, err := ah.userService.GetByID(context.Background(), *userID)
	if err != nil || currentUser.Role.Name != "db_admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	if err := ah.AdminService.OptimizeDatabase(context.Background()); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to optimize database"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Database optimized successfully"})
}
