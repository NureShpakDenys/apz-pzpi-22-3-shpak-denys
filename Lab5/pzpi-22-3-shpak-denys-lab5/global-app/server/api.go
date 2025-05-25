package main

import (
	"fmt"
	"net/http"
)

// APIError represents a generic error with an HTTP status code and a message.
type APIError struct {
	Status int
	Msg    string
}

// Error implements the error interface for APIError.
func (e *APIError) Error() string {
	return e.Msg
}

// NewAPIError constructs a new APIError with formatting.
func NewAPIError(status int, format string, args ...interface{}) error {
	return &APIError{
		Status: status,
		Msg:    fmt.Sprintf(format, args...),
	}
}

// NoFieldsFoundError represents a 404 error when a specific field or object is not found.
type NoFieldsFoundError struct {
	APIError
}

// NewNoFieldsFoundError returns a new NoFieldsFoundError for the given object name.
func NewNoFieldsFoundError(objName string) error {
	return &NoFieldsFoundError{
		APIError: APIError{
			Status: http.StatusNotFound,
			Msg:    fmt.Sprintf("no %s found", objName),
		},
	}
}
