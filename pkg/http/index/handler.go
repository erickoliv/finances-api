package index

import (
	"time"

	"github.com/gin-gonic/gin"
)

// Handler is the application root address. Can be used to check application status
func Handler(c *gin.Context) {
	current := time.Now().UTC()
	c.JSON(200, gin.H{
		"status": "ok",
		"utc":    current,
	})
}
