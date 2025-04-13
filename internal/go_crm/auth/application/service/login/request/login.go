package request

import (
	validator "go-api-docker/internal/common/validator"
)

type LoginRequestData struct {
	Email          string `json:"email" validate:"required,email"`
	Password       string `json:"password" validate:"required"`
	UserAgent      string
	AcceptLanguage string
}

type Login struct {
	requestData *LoginRequestData
	validator   *validator.Validator
}

// func CreatedLoginFromContext(c *gin.Context) *Login {
// 	var requestData LoginRequestData

// 	if err := c.ShouldBindJSON(&requestData); err != nil {
// 		userAgent := c.GetHeader("User-Agent")
// 		requestData.UserAgent = userAgent

// 		return &Login{
// 			requestData: requestData,
// 			err:         err,
// 			validator:   nil,
// 		}
// 	}

// 	validatorInstance := validator.NewValidator(request_factory.GetLanguage(c))
// 	return &Login{
// 		requestData: requestData,
// 		err:         nil,
// 		validator:   validatorInstance,
// 	}
// }

func CreatedLogin(requestData *LoginRequestData) *Login {
	validatorInstance := validator.NewValidator(requestData.AcceptLanguage)
	return &Login{
		requestData: requestData,
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

func (m *Login) GetValidationErrors() map[string]string {
	return m.validator.Validate(m.requestData)
}
