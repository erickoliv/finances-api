package auth

import (
	"net/http"

	"github.com/ericktm/olivsoft-golang-api/common"
	"github.com/ericktm/olivsoft-golang-api/model"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// Register creates a new user
func Register(c *gin.Context) {
	db := c.MustGet(common.DB).(*gorm.DB)
	user := model.User{}
	c.Bind(&user)

	credentials := Credentials{
		Username: user.Username,
		Password: user.Password,
	}
	credentials.Encrypt()

	user.Password = credentials.Password
	if err := db.Save(&user).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, common.ErrorMessage{
			Message: "registration error",
		})
	} else {
		c.JSON(http.StatusCreated, &user)
	}
}
