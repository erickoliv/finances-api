package auth

import (
	"net/http"
	"os"

	"github.com/ericktm/olivsoft-golang-api/api"
	"github.com/gin-gonic/gin"
)

// Middleware validates Authorization Headers
// TODO: Use JWT
func Middleware() gin.HandlerFunc {
	envToken := os.Getenv("APP_TOKEN")

	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")

		if token == "" {
			c.JSON(http.StatusUnauthorized, api.ErrorMessage{Message: "missing authentication token"})
			c.Abort()
			return
		}

		if token != envToken {
			c.JSON(http.StatusUnauthorized, api.ErrorMessage{Message: "invalid token"})
			c.Abort()
			return
		}

		c.Next()
	}
}
