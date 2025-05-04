package models // import "wayra/internal/core/domain/models"

import "gorm.io/gorm"

// User is a struct that represents the user model of the database
type User struct {
	// ID is the primary key of the user
	// Example: 1
	ID uint `gorm:"primaryKey;column:id"`

	// Name is the name of the user
	// Example: John Doe
	Name string `gorm:"size:255;not null;column:name; unique" json:"name"`

	// Password is the password of the user
	// Example: 111111
	Password string `gorm:"size:255;not null;column:password"`

	// RoleID is the foreign key of the role
	// Example: 1
	RoleID uint `gorm:"not null;column:role_id"`

	// Role is the role of the user
	Role Role `gorm:"foreignKey:RoleID" json:"role,omitempty"`

	// Companies is the companies of the user
	Companies []Company `gorm:"many2many:user_companies;constraint:OnDelete:CASCADE;" json:"companies,omitempty"`
}

// LoadRelations is an implementation of the interface LoadRelations for the User model
func (u *User) LoadRelations(db *gorm.DB) *gorm.DB {
	withRole := db.Preload("Role")
	withCompanies := withRole.Preload("Companies").Preload("Companies.Creator").Preload("Companies.Users")
	return withCompanies
}
