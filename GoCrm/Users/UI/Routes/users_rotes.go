package Routes

import (
	users_controllers "go-api-docker/GoCrm/Users/UI/Controllers"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(routerGroup *gin.RouterGroup) {
	routerGroup.POST("/user", users_controllers.CreateUser)
}
