package services

import (
	"context"
	"time"
)

type DBStatus struct {
	DatabaseSizeMB    float64
	ActiveConnections int
	LastBackupTime    time.Time
}

type AdminService interface {
	GetSystemConfig() map[string]interface{}
	UpdateSystemConfig(timeoutSeconds int, tokenTTLHrs int) error
	GetDBStatus(ctx context.Context) (*DBStatus, error)
	OptimizeDatabase(ctx context.Context) error
	ClearOldLogs(ctx context.Context, olderThanDays int) error
	BackupDatabase(ctx context.Context, userID uint, backupPath string) error
	RestoreDatabase(ctx context.Context, userID uint, backupPath string) error
}
