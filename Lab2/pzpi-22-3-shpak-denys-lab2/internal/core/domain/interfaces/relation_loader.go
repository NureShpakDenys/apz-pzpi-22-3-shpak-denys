// Package interfaces provides the interfaces for the domain layer
// It defines the interface for the RelationLoader
package interfaces // import "wayra/internal/core/domain/interfaces"

import "gorm.io/gorm"

// RelationLoader is an interface that defines the method to load relations
type RelationLoader interface {
	LoadRelations(db *gorm.DB) *gorm.DB
}
