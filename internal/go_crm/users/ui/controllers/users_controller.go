package controllers

import (
	"github.com/gin-gonic/gin"
)

func createUser(c *gin.Context) {
	c.JSON(200, gin.H{"data": "New user is created successfully!!!"})
}

func RegisterUserRoutes(router *gin.Engine) {
	group := router.Group("/api")
	{
		group.POST("/user", createUser)
		group.POST("/user-test", createUser)
	}
}
