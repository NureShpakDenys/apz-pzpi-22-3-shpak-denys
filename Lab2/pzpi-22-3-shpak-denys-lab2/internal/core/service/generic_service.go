package service // import "wayra/internal/core/service"

import (
	"context"
	"wayra/internal/core/port"
)

// GenericService is a generic service that implements the Service interface
type GenericService[T any] struct {
	Repository port.Repository[T] // Repository is the repository that the service will use
}

// NewGenericService creates a new GenericService
// repo: Repository that the service will use
// returns: a new GenericService
func NewGenericService[T any](repo port.Repository[T]) *GenericService[T] {
	return &GenericService[T]{Repository: repo}
}

// Create creates a new entity
// ctx: context
// entity: entity to create
// returns: error
func (s *GenericService[T]) Create(ctx context.Context, entity *T) error {
	return s.Repository.Add(ctx, entity)
}

// GetByID gets an entity by its ID
// ctx: context
// id: ID of the entity
// returns: the entity and an error
func (s *GenericService[T]) GetByID(ctx context.Context, id uint) (*T, error) {
	return s.Repository.GetByID(ctx, id)
}

// Where gets a list of entities that match the params
// ctx: context
// params: params to filter the entities
// returns: a list of entities and an error
func (s *GenericService[T]) Where(ctx context.Context, params *T) ([]T, error) {
	return s.Repository.Where(ctx, params)
}

// Update updates an entity
// ctx: context
// entity: entity to update
// returns: error
func (s *GenericService[T]) Update(ctx context.Context, entity *T) error {
	return s.Repository.Update(ctx, entity)
}

// Delete deletes an entity by its ID
// ctx: context
// id: ID of the entity
// returns: error
func (s *GenericService[T]) Delete(ctx context.Context, id uint) error {
	return s.Repository.Delete(ctx, id)
}

// SkipTake gets a list of entities skipping the first n elements and taking the next m elements
// ctx: context
// skip: number of elements to skip
// take: number of elements to take
// returns: a list of entities and an error
func (s *GenericService[T]) SkipTake(ctx context.Context, skip int, take int) (*[]T, error) {
	return s.Repository.SkipTake(ctx, skip, take)
}

// CountWhere counts the number of entities that match the params
// ctx: context
// params: params to filter the entities
// returns: the number of entities
func (s *GenericService[T]) CountWhere(ctx context.Context, params *T) int64 {
	return s.Repository.CountWhere(ctx, params)
}
