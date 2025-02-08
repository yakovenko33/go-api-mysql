package routes

import (
	users_controllers "go-api-docker/internal/go_crm/users/ui/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(routerGroup *gin.RouterGroup) {
	routerGroup.POST("/user", users_controllers.CreateUser)
}
