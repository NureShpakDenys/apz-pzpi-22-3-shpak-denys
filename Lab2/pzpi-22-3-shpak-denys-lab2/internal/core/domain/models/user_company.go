package models // import "wayra/internal/core/domain/models"

import "gorm.io/gorm"

// UserCompany is a struct that represent the user_company table into the database
type UserCompany struct {
	// ID is the id of the user_company
	// Example: 1
	ID uint `gorm:"primaryKey;column:id"`

	// UserID is the id of the user
	// Example: 1
	UserID uint `gorm:"not null;column:user_id"`

	// CompanyID is the id of the company
	// Example: 1
	CompanyID uint `gorm:"not null;column:company_id"`

	// Role is the role of the user in the company
	// Example: user
	Role string `gorm:"size:50;default:'user';column:role"`

	// User is the user of the user_company
	User User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;" json:"user,omitempty"`

	// Company is the company of the user_company
	Company Company `gorm:"foreignKey:CompanyID;constraint:OnDelete:CASCADE;" json:"company,omitempty"`
}

// LoadRelations is an implementation of the interface LoadRelations for the UserCompany and User struct
func (u *UserCompany) LoadRelations(db *gorm.DB) *gorm.DB {
	withUser := db.Preload("User")
	withCompany := withUser.Preload("Company")
	return withCompany
}
