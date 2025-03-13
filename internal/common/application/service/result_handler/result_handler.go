package result_handler

import (
	"errors"
	"net/http"
)

type ResultHandler[T any] struct {
	statusCode       int
	status           string
	errorString      string
	errorsValidation map[string]string
	result           any
}

func NewResultHandler[T any](statusCode int) *ResultHandler[T] {
	return &ResultHandler[T]{
		statusCode: statusCode,
	}
}

func (m *ResultHandler[T]) SetStatusCode(statusCode int) *ResultHandler[T] {
	m.statusCode = statusCode

	return m
}

func (m *ResultHandler[T]) GetStatusCode() int {
	return m.statusCode
}

func (m *ResultHandler[T]) SetError(errorString string) *ResultHandler[T] {
	m.errorString = errorString

	return m
}

func (m *ResultHandler[T]) GetError() string {
	return m.errorString
}

func (m *ResultHandler[T]) SetStatus(status string) *ResultHandler[T] {
	m.status = status

	return m
}

func (m *ResultHandler[T]) GetStatus() string {
	return m.status
}

func (m *ResultHandler[T]) SetErrorsValidation(errorsValidation map[string]string) *ResultHandler[T] {
	m.errorsValidation = errorsValidation

	return m
}

func (m *ResultHandler[T]) GetErrorsValidation() map[string]string {
	return m.errorsValidation
}

func (m *ResultHandler[T]) SetArrayResult(result []T) *ResultHandler[T] {
	m.result = result

	return m
}

func (m *ResultHandler[T]) SetSingleResult(result T) *ResultHandler[T] {
	m.result = result

	return m
}

func (m *ResultHandler[T]) GetResult() (T, bool) {
	res, ok := m.result.(T)
	return res, ok
}

const (
	StatusOk           = "ok"
	ValidateError      = "FieldValidateError"
	ServerError        = "ServerError"
	ForbiddenError     = "ForbiddenError"
	BusinessLogicError = "BusinessLogicError"
	DublicationError   = "DuplicationError"
)

type RequestInterface interface {
	GetError() error
	GetValidationErrors() map[string]string
}

func FactoryResultHandler[T any](request RequestInterface) (*ResultHandler[T], error) {
	errorsMap := request.GetValidationErrors()
	if len(errorsMap) > 0 {
		return NewResultHandler[T](http.StatusUnprocessableEntity).
			SetErrorsValidation(errorsMap).
			SetStatus(ValidateError), errors.New("has ErrorsValidation")
	}

	err := request.GetError()
	if err != nil {
		return NewResultHandler[T](http.StatusUnprocessableEntity).
			SetError(err.Error()).
			SetStatus(ServerError), errors.New("error mapping data")
	}

	return &ResultHandler[T]{}, nil
}
