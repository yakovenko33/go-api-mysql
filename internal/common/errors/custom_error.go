package errors

import (
	"fmt"
)

type CustomeError struct {
	Code    int
	Message string
	Status  string
}

func (e *CustomeError) Error() string {
	return fmt.Sprintf("code: %d, message: %s", e.Code, e.Message)
}

func NewCustomeError(code int, message string, status string) error {
	return &CustomeError{Code: code, Message: message, Status: status}
}
