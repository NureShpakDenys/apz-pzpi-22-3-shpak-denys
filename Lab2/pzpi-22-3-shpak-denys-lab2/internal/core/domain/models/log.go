package models

import "time"

type Log struct {
	// ID is the log identifier
	// Example: 1
	ID uint `gorm:"primaryKey;column:id"`

	// CreatedAt is the time when the log was created
	// Example: 2023-10-01T12:00:00Z
	CreatedAt time.Time `gorm:"autoCreateTime"`

	// UserID is the identifier of the user who performed the action it could be null
	// Example: 1
	UserID *uint `gorm:"column:user_id"`

	// ActionType is the type of action performed by the user
	// Example: "create", "update", "delete"
	ActionType string `gorm:"column:action_type"`

	// Description is a description of the action performed by the user
	// Example: "User created", "User updated", "User deleted"
	Description string `gorm:"column:description"`

	// Success indicates whether the action was successful or not
	// Example: true, false
	Success bool `gorm:"column:success"`
}
