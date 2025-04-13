package controllers

import (
	"errors"
	"net/http"
	"os"

	request_factory "go-api-docker/internal/common/application/service/request/request_factory"
	jwt_auth "go-api-docker/internal/common/security/auth/jwt_auth"
	middleware "go-api-docker/internal/common/server/middleware"
	response_helper "go-api-docker/internal/common/ui/response"
	login_handler "go-api-docker/internal/go_crm/auth/application/service/login"
	login_request "go-api-docker/internal/go_crm/auth/application/service/login/request"
	refresh_token "go-api-docker/internal/go_crm/auth/application/service/refresh_tokens/request"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	loginHandler *login_handler.LoginHandler
	jwtAuth      jwt_auth.JwtAuthManagerInterface
}

func NewAuthController(loginHandler *login_handler.LoginHandler, jwtAuth jwt_auth.JwtAuthManagerInterface) *AuthController {
	return &AuthController{
		loginHandler: loginHandler,
		jwtAuth:      jwtAuth,
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
		resultValue, _ := result.GetResult()
		c.SetCookie("refresh_token", resultValue.RefreshToken, int(resultValue.RefreshTokenExpiry), "/", os.Getenv("APP_DOMAIN"), true, true)
	}

	response_helper.Response(c, result)
}

func (m *AuthController) Loguout(c *gin.Context) {
	c.JSON(200, gin.H{"data": "Loguout"})
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
	refresh_token.CreateRefreshTokens(requestData)

	// tokens, err := m.jwtAuth.RefreshTokens(refreshToken, userData)
	// if err != nil {
	// 	c.JSON(401, gin.H{"error": err.Error()})
	// }
	setRefreshToken(c, &tokens)

	c.JSON(200, gin.H{"data": "refreshToken"})
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
		group.POST("/recoveryPassword", authController.RecoveryPassword)
	}
	groupSecure := router.Group("/api", middleware.CombinedMiddleware(jwtAuth)...)
	{
		groupSecure.POST("/loguout", authController.Loguout)
		groupSecure.POST("/refreshToken", authController.RefreshToken)
	}
}
