package response

import (
	result_handler "go-api-docker/internal/common/application/service/result_handler"

	"github.com/gin-gonic/gin"
)

func Response[T any](c *gin.Context, result *result_handler.ResultHandler[T]) {
	if len(result.GetErrorsValidation()) > 0 || result.GetError() != "" {
		c.JSON(result.GetStatusCode(), gin.H{
			"validation_errors": result.GetErrorsValidation(),
			"error":             result.GetError(),
			"result":            nil,
			"status":            result.GetStatus(),
		})
		return
	}

	resultValue, _ := result.GetResult()

	c.JSON(result.GetStatusCode(), gin.H{
		"data":              resultValue,
		"errors":            "",
		"validation_errors": nil,
		"status":            result,
	})
}
