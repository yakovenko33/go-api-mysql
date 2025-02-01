package Controllers

import (
	"github.com/gin-gonic/gin"
)

func CreateUser(c *gin.Context) {
	c.JSON(200, gin.H{"data": "New user is created successfully"})
}
