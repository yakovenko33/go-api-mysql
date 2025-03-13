package request

import (
	"github.com/gin-gonic/gin"

	validator "go-api-docker/internal/common/validator"
)

type requestData struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type Login struct {
	requestData requestData
	err         error
	validator   *validator.Validator
}

func CreatedLoginFromContext(c *gin.Context) *Login {
	var requestData requestData

	validatorInstance := validator.NewValidator("ru")
	if err := c.ShouldBindJSON(&requestData); err != nil {
		return &Login{
			requestData: requestData,
			err:         err,
			validator:   validatorInstance,
		}
	}

	return &Login{
		requestData: requestData,
		err:         nil,
		validator:   validatorInstance,
	}
}

func (m *Login) GetEmail() string {
	return m.requestData.Email
}

func (m *Login) GetPassword() string {
	return m.requestData.Password
}

func (m *Login) GetError() error {
	return m.err
}

func (m *Login) GetValidationErrors() map[string]string {
	return m.validator.Validate(m.requestData)
}
