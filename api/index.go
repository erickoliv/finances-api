package api

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// IndexHandler is the application root address. Can be used to check application status
func IndexHandler(app *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	}
}
