package url

import (
	"net/http"
	"os"

	"github.com/ericktm/olivsoft-golang-api/api"
	"github.com/ericktm/olivsoft-golang-api/database"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// PrepareRouter add description
func PrepareRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()
	r.GET("/", api.IndexHandler(db))

	r.Use(database.Middleware(db))
	r.Use(AuthMiddleware())

	return r
}

// AuthMiddleware validates Authorization Headers
func AuthMiddleware() gin.HandlerFunc {
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
