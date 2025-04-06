package controllers

import (
	"fmt"
	jwt_auth "go-api-docker/internal/common/security/auth/jwt_auth"
	middleware "go-api-docker/internal/common/server/middleware"
	response_helper "go-api-docker/internal/common/ui/response"
	login_handler "go-api-docker/internal/go_crm/auth/application/service/login"
	login_request "go-api-docker/internal/go_crm/auth/application/service/login/request"

	"net/http"

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
	loginRequest := login_request.CreatedLoginFromContext(c)

	result := m.loginHandler.Handle(loginRequest)
	if len(result.GetErrorsValidation()) == 0 && result.GetError() == "" {
		resultValue, _ := result.GetResult()
		c.SetCookie("refresh_token", resultValue.RefreshToken, int(resultValue.RefreshTokenExpiry), "/", "", false, true)
	}

	response_helper.Response(c, result)
}

func (m *AuthController) LoginTest(c *gin.Context) {
	fmt.Println("CreatedLoginFromContext")
	loginRequest := login_request.CreatedLoginFromContext(c)

	if loginRequest.GetError() != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": loginRequest.GetError().Error()})
		return
	}
	errors := loginRequest.GetValidationErrors()
	if len(errors) > 0 {
		c.JSON(400, gin.H{"errors": errors})
		return
	}

	c.JSON(200, gin.H{"data": "login1", "test": len(errors)})
}

func (m *AuthController) Loguout(c *gin.Context) {
	c.JSON(200, gin.H{"data": "Loguout"})
}

func (m *AuthController) RefreshToken(c *gin.Context) {
	cookie, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(401, gin.H{"error": "You does not have refresh_token"})
		return
	}

	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(401, gin.H{"error": "You does not have `userId`"})
		return
	}

	userData := &jwt_auth.UserData{
		UserId:    userId.(string),
		UserAgent: c.GetHeader("User-Agent"),
	}

	tokens, err := m.jwtAuth.GenerateTokens(cookie, userData)
	if err != nil {
		c.JSON(401, gin.H{"error": err.Error()})
	}

	c.SetCookie("refresh_token", tokens.RefreshToken, int(tokens.RefreshTokenExpiry), "/", "", false, true)

	c.JSON(200, gin.H{"data": "refreshToken"})
}

func (m *AuthController) RecoveryPassword(c *gin.Context) {
	c.JSON(200, gin.H{"data": "recoveryPassword"})
}

func RegisterAuthRoutes(router *gin.Engine, jwtAuth jwt_auth.JwtAuthManagerInterface, authController *AuthController) {
	group := router.Group("/api")
	{
		group.POST("/login", authController.Login)
		group.POST("/loginTest", authController.LoginTest)
		group.POST("/recoveryPassword", authController.RecoveryPassword)
	}
	groupSecure := router.Group("/api", middleware.CombinedMiddleware(jwtAuth)...)
	{
		groupSecure.POST("/loguout", authController.Loguout)
		groupSecure.POST("/refreshToken", authController.RefreshToken)
	}
}
