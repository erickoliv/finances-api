package auth

import (
	"net/http"

	"github.com/ericktm/olivsoft-golang-api/constants"
	"github.com/ericktm/olivsoft-golang-api/model"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// Register creates a new user
func Register(c *gin.Context) {
	db := c.MustGet(constants.DB).(*gorm.DB)
	user := model.User{}
	c.Bind(&user)
	if err := db.Save(&user).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
	} else {
		c.JSON(http.StatusCreated, &user)
	}
}
