package request

import (
	"github.com/gin-gonic/gin"

	request_factory "go-api-docker/internal/common/application/service/request/request_factory"
	validator "go-api-docker/internal/common/validator"
)

type requestData struct {
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required"`
	UserAgent string
}

type Login struct {
	requestData requestData
	err         error
	validator   *validator.Validator
}

func CreatedLoginFromContext(c *gin.Context) *Login {
	var requestData requestData

	if err := c.ShouldBindJSON(&requestData); err != nil {
		userAgent := c.GetHeader("User-Agent")
		requestData.UserAgent = userAgent

		return &Login{
			requestData: requestData,
			err:         err,
			validator:   nil,
		}
	}

	validatorInstance := validator.NewValidator(request_factory.GetLanguage(c))
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

func (m *Login) GetUserAgent() string {
	return m.requestData.UserAgent
}

func (m *Login) GetError() error {
	return m.err
}

func (m *Login) GetValidationErrors() map[string]string {
	return m.validator.Validate(m.requestData)
}
