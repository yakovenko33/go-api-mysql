package errors

import (
	"fmt"
)

type TokenInvalidError struct {
	Code    int
	Message string
}

func (e *TokenInvalidError) Error() string {
	return fmt.Sprintf("code: %d, message: %s", e.Code, e.Message)
}

func NewAuthError(code int, message string) error {
	return &TokenInvalidError{Code: code, Message: message}
}
