package admin

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
	"wayra/internal/adapter/httpserver/handlers"

	"github.com/gin-gonic/gin"
)

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
