package auth

import (
	"net/http"

	"github.com/ericktm/olivsoft-golang-api/common"
	"github.com/ericktm/olivsoft-golang-api/model"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// Login Route
func Login(c *gin.Context) {
	db := c.MustGet(common.DB).(*gorm.DB)
	credentials := Credentials{}
	user := model.User{}

	c.Bind(&credentials)
	credentials.Encrypt()

	if err := db.First(&user, "username = ? AND password = ?", credentials.Username, credentials.Password).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, common.ErrorMessage{
			Message: "login denied... add more details here",
		})
	} else {
		c.JSON(http.StatusOK, "return jwt token here")
	}
}
