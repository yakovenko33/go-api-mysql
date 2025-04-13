package refresh_tokens

import (
	result_handler "go-api-docker/internal/common/application/service/result_handler"
	jwt_auth "go-api-docker/internal/common/security/auth/jwt_auth"
	refresh_tokens "go-api-docker/internal/go_crm/auth/application/service/refresh_tokens/request"

	"go.uber.org/zap"
)

type RefreshTokensHandler struct {
	resultHandler *result_handler.ResultHandler[*jwt_auth.JwtTokens]
	jwtAuth       jwt_auth.JwtAuthManagerInterface
	logger        *zap.Logger
}

func NewLoginRefreshTokensHandler(
	jwtAuth jwt_auth.JwtAuthManagerInterface,
	logger *zap.Logger,
) *RefreshTokensHandler {
	return &RefreshTokensHandler{
		resultHandler: &result_handler.ResultHandler[*jwt_auth.JwtTokens]{},
		jwtAuth:       jwtAuth,
		logger:        logger,
	}
}

func (m *RefreshTokensHandler) Handle(request *refresh_tokens.RefreshTokens) *result_handler.ResultHandler[*jwt_auth.JwtTokens] {
	resultHandler, err := result_handler.FactoryResultHandler[*jwt_auth.JwtTokens](request)

	if err != nil {
		return resultHandler
	}
	m.resultHandler = resultHandler

	userData := &jwt_auth.UserData{
		UserId:    request.GetUserId(),
		UserAgent: request.GetUserAgent(),
	}
	tokens, err := m.jwtAuth.RefreshTokens(request.GetRefreshToken(), userData)

	return m.resultHandler.SetSingleResult(&tokens).SetStatus(result_handler.StatusOk).SetStatusCode(200)
}
