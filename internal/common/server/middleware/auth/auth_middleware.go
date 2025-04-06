package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	jwt_auth "go-api-docker/internal/common/security/auth/jwt_auth"
)

func AuthMiddleware(jwtAuth jwt_auth.JwtAuthManagerInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			notUnauthorized(c, "Missing token")
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			notUnauthorized(c, "Invalid token format")
			return
		}

		userId, err := jwtAuth.VerifyToken(parts[1])
		if err != nil || userId == "" {
			notUnauthorized(c, "Invalid or expired token")
			return
		}

		c.Set("userID", userId)
		c.Next()
	}
}

func notUnauthorized(c *gin.Context, errorText string) {
	c.JSON(http.StatusUnauthorized, gin.H{"error": errorText})
	c.Abort()
}
