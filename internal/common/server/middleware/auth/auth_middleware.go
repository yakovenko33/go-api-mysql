package auth

import (
	"github.com/gin-gonic/gin"

	jwt_auth "go-api-docker/internal/common/security/auth/jwt_auth"
)

func AuthMiddleware(jwt_auth *jwt_auth.JwtAuthManagerInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		//authHeader := c.GetHeader("Authorization")
		//reqToken := strings.Split(authHeader, " ")[1]
	}
}
