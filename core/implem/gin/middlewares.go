package coregin

import (
	"home-broker/core"

	"github.com/gin-gonic/gin"
)

// MiddlewareAPIError catch APIErrors and returns as JSON.
func MiddlewareAPIError() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		err := c.Errors.Last()
		if err == nil {
			return
		}
		apiError, ok := err.Err.(core.APIError)
		if !ok {
			apiError = core.NewAPIError("Internal server error.", 500)
		}
		c.JSON(apiError.StatusCode, gin.H{"error": gin.H{"message": apiError.Message}})
		c.Abort()
	}
}
