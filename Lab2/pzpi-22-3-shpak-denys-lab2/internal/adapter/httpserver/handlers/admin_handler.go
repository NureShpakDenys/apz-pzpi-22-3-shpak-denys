package handlers // import "wayra/internal/adapter/httpserver/handlers"

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"time"
	"wayra/internal/adapter/config"
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

// BackupDatabaseRequest is the request for the BackupDatabase endpoint
type BackupDatabaseRequest struct {
	// Backup path is the path where the backup will be stored
	BackupPath string `json:"backup_path"`
}

func encryptFile(filename string, key []byte) error {
	plaintext, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return err
	}

	ciphertext := aesGCM.Seal(nonce, nonce, plaintext, nil)

	return ioutil.WriteFile(filename, ciphertext, 0644)
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

	userID, err := getUserIDFromToken(c)
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
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	tables := []string{
		"roles",
		"users",
		"companies",
		"routes",
		"deliveries",
		"product_categories",
		"products",
		"waypoints",
		"sensor_data",
		"user_companies",
	}

	if err := os.MkdirAll(req.BackupPath, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create backup directory: " + err.Error()})
		ah.LogService.LogAction(userID, "backup_database", "Database backup failed: "+err.Error(), false)
		return
	}

	for _, table := range tables {
		filePath := fmt.Sprintf("%s/%s.csv", req.BackupPath, table)

		cmd := exec.Command(
			"psql",
			"-U", "postgres",
			"-d", "Wayra",
			"-h", "localhost",
			"-p", "5432",
			"-c", fmt.Sprintf(`\COPY %s TO '%s' WITH CSV HEADER`, table, filePath),
		)
		cmd.Env = append(os.Environ(), "PGPASSWORD="+ah.dbPassword)

		if err := cmd.Run(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": fmt.Sprintf("Error exporting table %s: %v", table, err),
			})
			ah.LogService.LogAction(userID, "backup_database", "Error exporting table: "+err.Error(), false)
			return
		}

		if err := encryptFile(filePath, ah.encryptionKey); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": fmt.Sprintf("Error encrypting file %s: %v", filePath, err),
			})
			ah.LogService.LogAction(userID, "backup_database", "Error encrypting file: "+err.Error(), false)
			return
		}
	}
	ah.db.Exec("INSERT INTO backup_logs (backup_time) VALUES (now())")
	c.JSON(http.StatusOK, "Backup created")
	ah.LogService.LogAction(userID, "backup_database", "Database backup created at "+req.BackupPath, true)
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

	userID, err := getUserIDFromToken(c)
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
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	tables := []string{
		"roles",
		"users",
		"companies",
		"routes",
		"deliveries",
		"product_categories",
		"products",
		"waypoints",
		"sensor_data",
		"user_companies",
	}

	tempPath := fmt.Sprintf("%s/temp", req.BackupPath)
	if err := os.MkdirAll(tempPath, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create temp directory: " + err.Error()})
		return
	}

	defer os.RemoveAll(tempPath)

	for _, table := range tables {
		encryptedFilePath := fmt.Sprintf("%s/%s.csv", req.BackupPath, table)
		tempFilePath := fmt.Sprintf("%s/%s.csv", tempPath, table)

		if err := decryptFileTo(encryptedFilePath, tempFilePath, ah.encryptionKey); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": fmt.Sprintf("Error decrypting file %s: %v", encryptedFilePath, err),
			})
			ah.LogService.LogAction(userID, "restore_database", "Database restore failed: "+err.Error(), false)
			return
		}

		truncateCmd := exec.Command(
			"psql",
			"-U", "postgres",
			"-d", "Wayra",
			"-h", "localhost",
			"-p", "5432",
			"-c", fmt.Sprintf(`TRUNCATE TABLE %s RESTART IDENTITY CASCADE;`, table),
		)
		truncateCmd.Env = append(os.Environ(), "PGPASSWORD="+ah.dbPassword)
		if output, err := truncateCmd.CombinedOutput(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": fmt.Sprintf("Error truncating table %s: %v, output: %s", table, err, string(output)),
			})
			ah.LogService.LogAction(userID, "restore_database", "Database restore failed: "+err.Error(), false)
			return
		}

		importCmd := exec.Command(
			"psql",
			"-U", "postgres",
			"-d", "Wayra",
			"-h", "localhost",
			"-p", "5432",
			"-c", fmt.Sprintf(`\COPY %s FROM '%s' WITH CSV HEADER`, table, tempFilePath),
		)
		importCmd.Env = append(os.Environ(), "PGPASSWORD="+ah.dbPassword)
		if output, err := importCmd.CombinedOutput(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": fmt.Sprintf("Error importing table %s: %v, output: %s", table, err, string(output)),
			})
			ah.LogService.LogAction(userID, "restore_database", "Database restore failed: "+err.Error(), false)
			return
		}
	}

	ah.LogService.LogAction(userID, "restore_database", "Database restored from "+req.BackupPath, true)
	c.JSON(http.StatusOK, "Database restored")
}

func decryptFileTo(encryptedPath, decryptedPath string, key []byte) error {
	ciphertext, err := ioutil.ReadFile(encryptedPath)
	if err != nil {
		return err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	nonceSize := aesGCM.NonceSize()
	if len(ciphertext) < nonceSize {
		return fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(decryptedPath, plaintext, 0644)
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

	userID, err := getUserIDFromToken(c)
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
	userID, err := getUserIDFromToken(c)
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
	userID, err := getUserIDFromToken(c)
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
	userID, err := getUserIDFromToken(c)
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

	userID, err := getUserIDFromToken(c)
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

// GetDBStatus godoc
// @Summary Get database status
// @Description Get database status
// @Tags admin
// @Accept json
// @Produce json
// @Security     BearerAuth
// @Router /admin/db-status [get]
func (ah *AdminHandler) GetDBStatus(c *gin.Context) {
	userID, err := getUserIDFromToken(c)
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
	userID, err := getUserIDFromToken(c)
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

	userID, err := getUserIDFromToken(c)
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
