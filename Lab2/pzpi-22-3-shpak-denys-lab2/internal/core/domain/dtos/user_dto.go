package dtos // import "wayra/internal/core/domain/dtos"

// UserDTO is a struct that represents the user data transfer object
type UserDTO struct {
	// ID is the user identifier
	// Example: 1
	ID uint `json:"id"`

	// Name is the user name
	// Example: John Doe
	Name string `json:"name"`
}
