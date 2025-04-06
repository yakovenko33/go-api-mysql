package middleware

import (
	"github.com/gin-gonic/gin"

	jwt_auth "go-api-docker/internal/common/security/auth/jwt_auth"
	auth_middleware "go-api-docker/internal/common/server/middleware/auth"
)

func CombinedMiddleware(jwtAuth jwt_auth.JwtAuthManagerInterface) []gin.HandlerFunc {
	return []gin.HandlerFunc{auth_middleware.AuthMiddleware(jwtAuth)}
}
