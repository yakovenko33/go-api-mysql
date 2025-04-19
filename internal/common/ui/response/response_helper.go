package response

import (
	result_handler "go-api-docker/internal/common/application/service/result_handler"

	"github.com/gin-gonic/gin"
)

func Response[T any](c *gin.Context, result *result_handler.ResultHandler[T]) {
	if len(result.GetErrorsValidation()) > 0 || result.GetError() != "" {
		c.JSON(result.GetStatusCode(), gin.H{
			"data":              nil,
			"validation_errors": result.GetErrorsValidation(),
			"error":             result.GetError(),
			"status":            result.GetStatus(),
		})
		return
	}

	resultValue, _ := result.GetResult()

	c.JSON(result.GetStatusCode(), gin.H{
		"data":              resultValue,
		"error":             "",
		"validation_errors": nil,
		"status":            result.GetStatus(),
	})
}

func ResponseServerError(c *gin.Context, statusCode int, err error) {
	c.JSON(statusCode, gin.H{
		"data":              nil,
		"error":             err.Error(),
		"validation_errors": nil,
		"status":            result_handler.ServerError,
	})
}

func ResponseServerSuccessful[T any](c *gin.Context, data T, statusCode int) {
	c.JSON(statusCode, gin.H{
		"data":              data,
		"error":             nil,
		"validation_errors": nil,
		"status":            result_handler.StatusOk,
	})
}
