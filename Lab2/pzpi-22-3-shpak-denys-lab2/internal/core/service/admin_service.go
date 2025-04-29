// internal/port/services/admin_service.go
package service

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"time"

	"wayra/internal/adapter/config"
	"wayra/internal/core/port/services"

	"gorm.io/gorm"
)

type AdminService struct {
	cfg           *config.Config
	db            *gorm.DB
	EncryptionKey []byte
}

func NewAdminService(cfg *config.Config, db *gorm.DB, encKey []byte) services.AdminService {
	return &AdminService{cfg: cfg, db: db, EncryptionKey: encKey}
}

var backupTables = []string{
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

func (s *AdminService) BackupDatabase(ctx context.Context, userID uint, backupPath string) error {
	if err := os.MkdirAll(backupPath, 0755); err != nil {
		return fmt.Errorf("failed to create backup directory: %w", err)
	}

	for _, table := range backupTables {
		filePath := fmt.Sprintf("%s/%s.csv", backupPath, table)

		if err := s.exportTable(table, filePath); err != nil {
			return fmt.Errorf("error exporting table %s: %w", table, err)
		}

		if err := encryptFile(filePath, s.EncryptionKey); err != nil {
			return fmt.Errorf("error encrypting file %s: %w", filePath, err)
		}
	}

	s.db.Exec("INSERT INTO backup_logs (backup_time) VALUES (now())")
	return nil
}

func (s *AdminService) RestoreDatabase(ctx context.Context, userID uint, backupPath string) error {
	tempPath := fmt.Sprintf("%s/temp", backupPath)
	if err := os.MkdirAll(tempPath, 0755); err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tempPath)

	for _, table := range backupTables {
		encPath := fmt.Sprintf("%s/%s.csv", backupPath, table)
		tmpPath := fmt.Sprintf("%s/%s.csv", tempPath, table)

		if err := decryptFileTo(encPath, tmpPath, s.EncryptionKey); err != nil {
			return fmt.Errorf("error decrypting file %s: %w", encPath, err)
		}

		if err := s.truncateTable(table); err != nil {
			return fmt.Errorf("error truncating table %s: %w", table, err)
		}

		if err := s.importTable(table, tmpPath); err != nil {
			return fmt.Errorf("error importing table %s: %w", table, err)
		}
	}

	return nil
}

func (s *AdminService) GetSystemConfig() map[string]interface{} {
	return map[string]interface{}{
		"http_timeout_seconds":  int(s.cfg.Http.Timeout.Seconds()),
		"auth_token_ttl_hours":  int(s.cfg.AuthConfig.TokenExpiry.Hours()),
		"encryption_key_exists": s.cfg.EncryptionKey != "",
	}
}

func (s *AdminService) UpdateSystemConfig(timeoutSeconds int, tokenTTLHrs int) error {
	if timeoutSeconds <= 0 || tokenTTLHrs <= 0 {
		return errors.New("invalid timeout or TTL")
	}

	s.cfg.Http.Timeout = time.Duration(timeoutSeconds) * time.Second
	s.cfg.AuthConfig.TokenExpiry = time.Duration(tokenTTLHrs) * time.Hour
	return nil
}

func (s *AdminService) GetDBStatus(ctx context.Context) (*services.DBStatus, error) {
	var size float64
	err := s.db.Raw(`
		SELECT pg_database_size(current_database()) / 1024 / 1024 AS size_mb
	`).Scan(&size).Error
	if err != nil {
		return nil, err
	}

	var activeConnections int
	err = s.db.Raw(`
		SELECT count(*) FROM pg_stat_activity
	`).Scan(&activeConnections).Error
	if err != nil {
		return nil, err
	}

	var lastBackup time.Time
	err = s.db.Raw(`
		 SELECT backup_time FROM backup_logs ORDER BY backup_time DESC LIMIT 1
	`).Scan(&lastBackup).Error
	if err != nil || lastBackup.IsZero() {
		lastBackup = time.Time{}
	}
	return &services.DBStatus{
		DatabaseSizeMB:    size,
		ActiveConnections: activeConnections,
		LastBackupTime:    lastBackup,
	}, nil
}

func (s *AdminService) OptimizeDatabase(ctx context.Context) error {
	return s.db.Exec("VACUUM FULL").Error
}

func (s *AdminService) ClearOldLogs(ctx context.Context, olderThanDays int) error {
	if olderThanDays <= 0 {
		return errors.New("invalid number of days")
	}
	cutoff := time.Now().AddDate(0, 0, -olderThanDays)
	return s.db.Exec("DELETE FROM logs WHERE created_at < ?", cutoff).Error
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

func (s *AdminService) exportTable(table, filePath string) error {
	cmd := exec.Command(
		"psql",
		"-U", "postgres",
		"-d", "Wayra",
		"-h", "localhost",
		"-p", "5432",
		"-c", fmt.Sprintf(`\COPY %s TO '%s' WITH CSV HEADER`, table, filePath),
	)
	cmd.Env = append(os.Environ(), "PGPASSWORD="+s.cfg.DBPassword)
	return cmd.Run()
}

func (s *AdminService) truncateTable(table string) error {
	cmd := exec.Command(
		"psql",
		"-U", "postgres",
		"-d", "Wayra",
		"-h", "localhost",
		"-p", "5432",
		"-c", fmt.Sprintf(`TRUNCATE TABLE %s RESTART IDENTITY CASCADE;`, table),
	)
	cmd.Env = append(os.Environ(), "PGPASSWORD="+s.cfg.DBPassword)
	_, err := cmd.CombinedOutput()
	return err
}

func (s *AdminService) importTable(table, filePath string) error {
	cmd := exec.Command(
		"psql",
		"-U", "postgres",
		"-d", "Wayra",
		"-h", "localhost",
		"-p", "5432",
		"-c", fmt.Sprintf(`\COPY %s FROM '%s' WITH CSV HEADER`, table, filePath),
	)
	cmd.Env = append(os.Environ(), "PGPASSWORD="+s.cfg.DBPassword)
	_, err := cmd.CombinedOutput()
	return err
}
