package auth

import (
	"fmt"
	"net/http"
	"os"

	"github.com/ericktm/olivsoft-golang-api/common"
	"github.com/ericktm/olivsoft-golang-api/pkg/domain"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// Register creates a new user
func Register(c *gin.Context) {
	db := c.MustGet(common.DB).(*gorm.DB)
	salt := os.Getenv(common.AppToken)
	user := domain.User{}

	if err := c.Bind(&user); err != nil {
		msg := fmt.Sprintf("invalid payload: %v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, common.ErrorMessage{
			Message: msg,
		})
		return
	}

	credentials := Credentials{
		Username: user.Username,
		Password: user.Password,
	}
	credentials.Encrypt(salt)

	user.Password = credentials.Password
	if err := db.Save(&user).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, common.ErrorMessage{
			Message: "registration error",
		})
		return
	}

	c.JSON(http.StatusCreated, &user)
}
