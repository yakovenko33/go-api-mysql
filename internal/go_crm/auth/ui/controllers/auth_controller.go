package controllers

import (
	"fmt"
	response_helper "go-api-docker/internal/common/ui/response"
	login_handler "go-api-docker/internal/go_crm/auth/application/service/login"
	login_request "go-api-docker/internal/go_crm/auth/application/service/login/request"

	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	loginHandler *login_handler.LoginHandler
}

func NewAuthController(loginHandler *login_handler.LoginHandler) *AuthController {
	return &AuthController{
		loginHandler: loginHandler,
	}
}

func (m *AuthController) Login(c *gin.Context) {
	loginRequest := login_request.CreatedLoginFromContext(c)

	result := m.loginHandler.Handle(loginRequest)

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
	c.JSON(200, gin.H{"data": "refreshToken"})
}

func (m *AuthController) RecoveryPassword(c *gin.Context) {
	c.JSON(200, gin.H{"data": "recoveryPassword"})
}

func RegisterAuthRoutes(router *gin.Engine, authController *AuthController) {
	group := router.Group("/api")
	{
		group.POST("/login", authController.Login)
		group.POST("/loginTest", authController.LoginTest)
		group.POST("/loguout", authController.Loguout)
		group.POST("/refreshToken", authController.RefreshToken)
		group.POST("/recoveryPassword", authController.RecoveryPassword)
	}
}
