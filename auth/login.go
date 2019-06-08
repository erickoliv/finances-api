package auth

import (
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
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
		ttl := time.Now().Add(60 * time.Minute)
		cookie := &Jwt{
			Username: credentials.Username,
			StandardClaims: jwt.StandardClaims{
				// In JWT, the expiry time is expressed as unix milliseconds
				ExpiresAt: ttl.Unix(),
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, cookie)

		key := []byte(os.Getenv(common.AppToken))
		if str, err := token.SignedString(key); err != nil {
			// If there is an error in creating the JWT return an internal server error
			c.AbortWithStatusJSON(http.StatusUnauthorized, common.ErrorMessage{
				Message: err.Error(),
			})
		} else {
			expiration := int(ttl.Sub(time.Now()).Seconds())
			c.SetCookie(common.AuthCookie, str, expiration, "/", "", false, true)
			c.JSON(http.StatusOK, ":)")
		}

	}
}
