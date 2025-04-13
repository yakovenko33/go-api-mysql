package login

import (
	"errors"
	result_handler "go-api-docker/internal/common/application/service/result_handler"
	custome_error "go-api-docker/internal/common/errors"
	auth_error "go-api-docker/internal/common/security/auth/errors"
	jwt_auth "go-api-docker/internal/common/security/auth/jwt_auth"
	password "go-api-docker/internal/common/security/password"
	login_request "go-api-docker/internal/go_crm/auth/application/service/login/request"
	users_entities "go-api-docker/internal/go_crm/users/domains/entities"
	users_repository "go-api-docker/internal/go_crm/users/infrastructure/repositories/users_repository"

	"fmt"

	"go.uber.org/zap"
)

type LoginHandler struct {
	resultHandler   *result_handler.ResultHandler[*jwt_auth.JwtTokens]
	usersRepository users_repository.UsersRepositoryInterface
	jwtAuth         jwt_auth.JwtAuthManagerInterface
	logger          *zap.Logger
}

func NewLoginHandler(
	usersRepository users_repository.UsersRepositoryInterface,
	jwtAuth jwt_auth.JwtAuthManagerInterface,
	logger *zap.Logger,
) *LoginHandler {
	return &LoginHandler{
		resultHandler:   &result_handler.ResultHandler[*jwt_auth.JwtTokens]{},
		usersRepository: usersRepository,
		jwtAuth:         jwtAuth,
		logger:          logger,
	}
}

func (m *LoginHandler) Handle(request *login_request.Login) *result_handler.ResultHandler[*jwt_auth.JwtTokens] {
	resultHandler, err := result_handler.FactoryResultHandler[*jwt_auth.JwtTokens](request)

	if err != nil {
		return resultHandler
	}
	m.resultHandler = resultHandler

	loginResult, err := m.login(request)
	var customErr *custome_error.CustomeError
	if err != nil && errors.As(err, &customErr) {
		return m.resultHandler.SetStatus(customErr.Status).SetError(err.Error()).SetStatusCode(customErr.Code)
	}

	var authError *auth_error.TokenInvalidError
	if err != nil && errors.As(err, &authError) {
		return m.resultHandler.SetStatus(result_handler.ServerError).SetError(err.Error()).SetStatusCode(400)
	}

	return m.resultHandler.SetSingleResult(loginResult).SetStatus(result_handler.StatusOk).SetStatusCode(200)
}

func (m *LoginHandler) login(request *login_request.Login) (*jwt_auth.JwtTokens, error) {
	user, err := m.findUserByEmail(request.GetEmail())
	if user == nil {
		return &jwt_auth.JwtTokens{}, err
	}

	if !password.CheckPasswordHash(request.GetPassword(), user.Password) {
		return &jwt_auth.JwtTokens{}, err
	}

	return m.generateTokens(user, request)
}

func (m *LoginHandler) findUserByEmail(email string) (*users_entities.User, error) {
	user, err := m.usersRepository.FindUserByEmail(email)

	if err != nil {
		m.logger.Error(fmt.Sprintf("Error for usersRepository.FindUserByEmail %s", err))
		return &users_entities.User{}, newCustomeError(500, "Problem on server. Try next time.", result_handler.ServerError)
	}
	if user == nil {
		return &users_entities.User{}, newCustomeError(404, fmt.Sprintf("User by email %s not found.", email), result_handler.BusinessLogicError)
	}

	return user, nil
}

func (m *LoginHandler) generateTokens(user *users_entities.User, request *login_request.Login) (*jwt_auth.JwtTokens, error) {
	userData := &jwt_auth.UserData{
		UserId:    user.ID.String(),
		UserAgent: request.GetUserAgent(),
	}

	tokens, err := m.jwtAuth.GenerateTokens(userData)
	if err == nil {
		return &jwt_auth.JwtTokens{}, err
	}

	return &tokens, nil
}

func newCustomeError(code int, message string, statusCode string) error {
	return custome_error.NewCustomeError(code, message, statusCode)
}
