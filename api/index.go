package api

import (
	"time"

	"github.com/gin-gonic/gin"
)

// IndexHandler is the application root address. Can be used to check application status
func IndexHandler(c *gin.Context) {
	current := time.Now().UTC()
	c.JSON(200, gin.H{
		"status": "ok",
		"utc":    current,
	})
}
