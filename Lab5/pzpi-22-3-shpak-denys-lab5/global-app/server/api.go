package main

import (
	"fmt"
	"net/http"
)

type APIError struct {
	Status int
	Msg    string
}

func (e *APIError) Error() string {
	return e.Msg
}

func NewAPIError(status int, format string, args ...interface{}) error {
	return &APIError{
		Status: status,
		Msg:    fmt.Sprintf(format, args...),
	}
}

type NoFieldsFoundError struct {
	APIError
}

func NewNoFieldsFoundError(objName string) error {
	return &NoFieldsFoundError{
		APIError: APIError{
			Status: http.StatusNotFound,
			Msg:    fmt.Sprintf("no %s found", objName),
		},
	}
}
