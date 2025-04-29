package models // import "wayra/internal/core/domain/models"

// Role is a struct that represent the role model
type Role struct {
	// ID is the role identifier
	// Example: 1
	ID uint `gorm:"primaryKey;column:id"`

	// Name is the role name
	// Example: admin
	Name string `gorm:"size:255;not null;column:name" json:"name"`

	// Description is the role description
	Users []User `gorm:"foreignKey:RoleID" json:"users,omitempty"`
}
