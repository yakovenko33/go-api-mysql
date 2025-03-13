package controllers

import (
	"fmt"
	response_helper "go-api-docker/internal/common/ui/response"
	login_handler "go-api-docker/internal/go_crm/auth/application/service/login"
	login_request "go-api-docker/internal/go_crm/auth/application/service/login/request"

	"net/http"

	"github.com/gin-gonic/gin"
)

func loginT(c *gin.Context) {
	loginRequest := login_request.CreatedLoginFromContext(c)

	handler := login_handler.NewLoginHandler()
	result := handler.Handle(loginRequest)

	response_helper.Response(c, result)
}

func login(c *gin.Context) {
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

func loguout(c *gin.Context) {
	c.JSON(200, gin.H{"data": "loguout"})
}

func refreshToken(c *gin.Context) {
	c.JSON(200, gin.H{"data": "refreshToken"})
}

func recoveryPassword(c *gin.Context) {
	c.JSON(200, gin.H{"data": "recoveryPassword"})
}

func RegisterAuthRoutes(router *gin.Engine) {
	group := router.Group("/api")
	{
		group.POST("/login", login)
		group.POST("/loginT", loginT)
		group.POST("/loguout", loguout)
		group.POST("/refreshToken", refreshToken)
		group.POST("/recoveryPassword", recoveryPassword)
	}
}
