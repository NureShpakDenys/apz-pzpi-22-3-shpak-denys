// Package port provides the repository interface for the core domain.
package port // import "wayra/internal/core/port"

import (
	"context"

	"gorm.io/gorm"
)

// Repository is the interface that wraps the basic CRUD operations.
type Repository[T any] interface {
	Add(ctx context.Context, entity *T) error
	GetByID(ctx context.Context, id uint) (*T, error)
	GetAll(ctx context.Context) ([]T, error)
	Where(ctx context.Context, params interface{}, args ...interface{}) ([]T, error)
	Update(ctx context.Context, entity *T) error
	Delete(ctx context.Context, id uint) error
	SkipTake(ctx context.Context, skip int, take int) (*[]T, error)
	CountWhere(ctx context.Context, params *T) int64
	DB() *gorm.DB
}
