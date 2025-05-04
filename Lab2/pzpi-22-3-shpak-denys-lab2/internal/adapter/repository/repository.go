package repository // import "wayra/internal/adapter/repository"

import (
	"context"
	"errors"
	"wayra/internal/core/domain/interfaces"
	"wayra/internal/core/port"

	"gorm.io/gorm"
)

// GenericRepository is a generic repository that implements the Repository interface
type GenericRepository[T any] struct {
	db *gorm.DB // db is the database connection
}

// NewRepository creates a new GenericRepository
// db: database connection
// returns: *GenericRepository
func NewRepository[T any](db *gorm.DB) port.Repository[T] {
	return &GenericRepository[T]{db: db}
}

func (r *GenericRepository[T]) DB() *gorm.DB {
	return r.db
}

// Add adds a new entity to the database
// ctx: context
// entity: entity to add
// returns: error
func (r *GenericRepository[T]) Add(ctx context.Context, entity *T) error {
	err := r.db.WithContext(ctx).Create(entity).Error
	if err != nil {
		return err
	}

	if loader, ok := any(entity).(interfaces.RelationLoader); ok {
		return loader.LoadRelations(r.db.WithContext(ctx)).First(entity).Error
	}

	return nil
}

// GetAll returns all entities from the database
// ctx: context
// id: id of the entity
// returns: []*T, error
func (r *GenericRepository[T]) GetByID(ctx context.Context, id uint) (*T, error) {
	var entity T
	query := r.db.WithContext(ctx).Model(&entity).Where("id = ?", id)

	if loader, ok := any(&entity).(interfaces.RelationLoader); ok {
		query = loader.LoadRelations(query)
	}

	err := query.First(&entity).Error
	if err != nil {
		return nil, err
	}

	return &entity, nil
}

// GetAll returns all entities from the database
// It also loads the relations of the entities
// It is possible to provide as object to filter the entities and string of type "where clause"
// ctx: context
// params: parameters to filter the entities
// args: arguments to filter the entities
// returns: []*T, error
func (r *GenericRepository[T]) Where(ctx context.Context, params interface{}, args ...interface{}) ([]T, error) {
	var entities []T
	query := r.db.WithContext(ctx)

	switch p := params.(type) {
	case *T:
		query = query.Where(p)
	case string:
		query = query.Where(p, args...)
	default:
		return nil, errors.New("unsupported parameter type for Where")
	}

	err := query.Find(&entities).Error
	if err != nil {
		return nil, err
	}

	for i := range entities {
		if loader, ok := any(&entities[i]).(interfaces.RelationLoader); ok {
			err := loader.LoadRelations(r.db.WithContext(ctx)).First(&entities[i]).Error
			if err != nil {
				return nil, err
			}
		}
	}

	return entities, nil
}

// Update updates an entity in the database
// ctx: context
// entity: entity to update
// returns: error
func (r *GenericRepository[T]) Update(ctx context.Context, entity *T) error {
	err := r.db.WithContext(ctx).Model(entity).Updates(entity).Error
	if err != nil {
		return err
	}
	if loader, ok := any(entity).(interfaces.RelationLoader); ok {
		return loader.LoadRelations(r.db.WithContext(ctx)).First(entity).Error
	}
	return nil
}

// Delete deletes an entity from the database
// ctx: context
// id: id of the entity
// returns: error
func (r *GenericRepository[T]) Delete(ctx context.Context, id uint) error {
	var entity T
	err := r.db.WithContext(ctx).Where("id = ?", id).Delete(&entity).Error
	if err != nil {
		return err
	}

	return nil
}

// SkipTake returns a slice of entities from the database
// ctx: context
// skip: number of entities to skip
// take: number of entities to take
// returns: []*T, error
func (r *GenericRepository[T]) SkipTake(ctx context.Context, skip int, take int) (*[]T, error) {
	var entities []T
	query := r.db.WithContext(ctx).Offset(skip).Limit(take)

	err := query.Find(&entities).Error
	if err != nil {
		return nil, err
	}
	return &entities, nil
}

// CountWhere returns the number of entities that match the query
// ctx: context
// params: parameters to filter the entities
// returns: int64
func (r *GenericRepository[T]) CountWhere(ctx context.Context, params *T) int64 {
	var count int64
	query := r.db.WithContext(ctx).Model(new(T)).Where(params)

	if loader, ok := any(params).(interfaces.RelationLoader); ok {
		query = loader.LoadRelations(query)
	}

	query.Count(&count)
	return count
}

// GetAll returns all entities from the database
// ctx: context
// returns: []*T, error
func (r *GenericRepository[T]) GetAll(ctx context.Context) ([]T, error) {
	var entities []T
	query := r.db.WithContext(ctx)

	err := query.Find(&entities).Error
	if err != nil {
		return nil, err
	}

	for i := range entities {
		if loader, ok := any(&entities[i]).(interfaces.RelationLoader); ok {
			err := loader.LoadRelations(r.db.WithContext(ctx)).First(&entities[i]).Error
			if err != nil {
				return nil, err
			}
		}
	}

	return entities, nil
}
