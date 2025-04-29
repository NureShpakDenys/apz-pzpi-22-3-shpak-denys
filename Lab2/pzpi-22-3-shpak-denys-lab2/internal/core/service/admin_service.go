// internal/port/services/admin_service.go
package service

import (
	"context"
	"errors"
	"time"

	"wayra/internal/adapter/config"
	"wayra/internal/core/port/services"

	"gorm.io/gorm"
)

type AdminService struct {
	cfg *config.Config
	db  *gorm.DB
}

func NewAdminService(cfg *config.Config, db *gorm.DB) services.AdminService {
	return &AdminService{cfg: cfg, db: db}
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
