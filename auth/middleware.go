package auth

import (
	"net/http"
	"os"

	"github.com/ericktm/olivsoft-golang-api/common"
	"github.com/gin-gonic/gin"
)

// Middleware to validate Authorization Headers
// TODO: Use JWT
func Middleware() gin.HandlerFunc {
	envToken := os.Getenv("APP_TOKEN")

	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")

		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, common.ErrorMessage{Message: "missing authentication token"})
			return
		}

		if token != envToken {
			c.AbortWithStatusJSON(http.StatusUnauthorized, common.ErrorMessage{Message: "invalid token"})
			return
		}

		c.Next()
	}
}
