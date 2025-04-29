// Package handlers contains the http handlers for the application
//
// This package is part of the adapter layer
// You can get handler instances by calling the New*Handler functions
// The handlers are responsible for handling the http requests and responses
package handlers // import "wayra/internal/adapter/httpserver/handlers"

// The role of user is one of these constants
const (
	AdminRole = iota + 1
	UserRole
	ManagerRole
)

// Role is a type for the user role in company
type Role string

// The role of user in company is one of these constants
const (
	RoleUser        Role = "user"
	RoleAdmin       Role = "admin"
	RoleManager     Role = "manager"
	RoleDBAdmin     Role = "db_admin"
	RoleSystemAdmin Role = "system_admin"
)

// The status of delivery is one of these constants
const (
	NotStarted = "not_started"
	InProgress = "in_progress"
	Completed  = "completed"
)
