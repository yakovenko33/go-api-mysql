package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	request_factory "go-api-docker/internal/common/application/service/request/request_factory"
	jwt_auth "go-api-docker/internal/common/security/auth/jwt_auth"
	middleware "go-api-docker/internal/common/server/middleware"
	response_helper "go-api-docker/internal/common/ui/response"
	login_handler "go-api-docker/internal/go_crm/auth/application/service/login"
	login_request "go-api-docker/internal/go_crm/auth/application/service/login/request"
	refresh_token_handler "go-api-docker/internal/go_crm/auth/application/service/refresh_tokens"
	refresh_token "go-api-docker/internal/go_crm/auth/application/service/refresh_tokens/request"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AuthController struct {
	loginHandler        *login_handler.LoginHandler
	jwtAuth             jwt_auth.JwtAuthManagerInterface
	refreshTokenHandler *refresh_token_handler.RefreshTokensHandler
	logger              *zap.Logger
}

func NewAuthController(
	loginHandler *login_handler.LoginHandler,
	jwtAuth jwt_auth.JwtAuthManagerInterface,
	refreshTokenHandler *refresh_token_handler.RefreshTokensHandler,
	logger *zap.Logger,
) *AuthController {
	return &AuthController{
		loginHandler:        loginHandler,
		jwtAuth:             jwtAuth,
		refreshTokenHandler: refreshTokenHandler,
		logger:              logger,
	}
}

func (m *AuthController) Login(c *gin.Context) {
	var requestData login_request.LoginRequestData
	if err := c.ShouldBindJSON(&requestData); err != nil {
		response_helper.ResponseServerError(c, http.StatusBadRequest, err)
		return
	}

	requestData.UserAgent = c.GetHeader("User-Agent")
	requestData.AcceptLanguage = request_factory.GetLanguage(c)

	loginRequest := login_request.CreatedLogin(&requestData)

	result := m.loginHandler.Handle(loginRequest)
	if len(result.GetErrorsValidation()) == 0 && result.GetError() == "" {
		tokens, _ := result.GetResult()
		setRefreshToken(c, tokens)
	}

	response_helper.Response(c, result)
}

func (m *AuthController) Loguout(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		response_helper.ResponseServerError(c, http.StatusUnauthorized, errors.New("you does not have refresh_token"))
		return
	}
	err = m.jwtAuth.AddToBlackList(refreshToken)
	if err != nil {
		m.logger.Error(fmt.Sprintf("Loguout %s", err))
		response_helper.ResponseServerError(c, http.StatusUnauthorized, errors.New("Loguout is not not successful"))
		return
	}

	response_helper.ResponseServerSuccessful(c, "Loguout is successful", http.StatusOK)
}

func (m *AuthController) RefreshToken(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		response_helper.ResponseServerError(c, http.StatusUnauthorized, errors.New("you does not have refresh_token"))
		return
	}

	userId, exists := c.Get("userId")
	if !exists {
		response_helper.ResponseServerError(c, http.StatusUnauthorized, errors.New("you does not have `userId`"))
		return
	}

	requestData := &refresh_token.RefreshTokensRequestData{
		UserId:       userId.(string),
		UserAgent:    c.GetHeader("User-Agent"),
		RefreshToken: refreshToken,
	}
	request := refresh_token.CreateRefreshTokens(requestData)
	result := m.refreshTokenHandler.Handle(request)

	tokens, _ := result.GetResult()
	setRefreshToken(c, tokens)

	response_helper.Response(c, result)
}

func setRefreshToken(c *gin.Context, tokens *jwt_auth.JwtTokens) {
	c.SetCookie("refresh_token", tokens.RefreshToken, int(tokens.RefreshTokenExpiry), "/", os.Getenv("APP_DOMAIN"), true, true)
}

func (m *AuthController) RecoveryPassword(c *gin.Context) {
	c.JSON(200, gin.H{"data": "recoveryPassword"})
}

func RegisterAuthRoutes(router *gin.Engine, jwtAuth jwt_auth.JwtAuthManagerInterface, authController *AuthController) {
	group := router.Group("/api")
	{
		group.POST("/login", authController.Login)
		group.POST("/recovery-password", authController.RecoveryPassword)
		group.POST("/refresh-token", authController.RefreshToken)
	}
	groupSecure := router.Group("/api", middleware.CombinedMiddleware(jwtAuth)...)
	{
		groupSecure.POST("/loguout", authController.Loguout)
	}
}
