package admin

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"
	"wayra/internal/adapter/httpserver/handlers"
	"wayra/internal/core/domain/models"

	"github.com/gin-gonic/gin"
)

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

// SendConfigToGlobalStorage godoc
// @Summary Send config to global storage
// @Description Sends current config to external storage
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Router /admin/send-config [post]
func (ah *AdminHandler) SendConfigToGlobalStorage(c *gin.Context) {
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

	type credsPayload struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Config   any    `json:"config"`
	}

	payload := credsPayload{
		Username: currentUser.Name,
		Password: currentUser.Password,
		Config:   ah.Cfg,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal config"})
		return
	}

	resp, err := http.Post("http://localhost:8080/set-creds", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send config"})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send config", "details": string(body)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Config sent successfully"})
}
