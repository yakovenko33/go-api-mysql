package login

import (
	result_handler "go-api-docker/internal/common/application/service/result_handler"
	jwt_auth "go-api-docker/internal/common/security/auth/jwt_auth"
	login_request "go-api-docker/internal/go_crm/auth/application/service/login/request"
)

type LoginHandler struct {
	resultHandler *result_handler.ResultHandler[jwt_auth.JwtTokens]
}

func NewLoginHandler() *LoginHandler {
	return &LoginHandler{
		resultHandler: &result_handler.ResultHandler[jwt_auth.JwtTokens]{},
	}
}

func (m *LoginHandler) Handle(request *login_request.Login) *result_handler.ResultHandler[jwt_auth.JwtTokens] {
	resultHandler, err := result_handler.FactoryResultHandler[jwt_auth.JwtTokens](request)

	if err != nil {
		return resultHandler
	}
	m.resultHandler = resultHandler

	return m.resultHandler
}
