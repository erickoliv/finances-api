package auth

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/ericktm/olivsoft-golang-api/common"
	"github.com/ericktm/olivsoft-golang-api/pkg/domain"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// Login Route
func Login(c *gin.Context) {
	db := c.MustGet(common.DB).(*gorm.DB)
	salt := os.Getenv(common.AppToken)
	credentials := Credentials{}
	user := domain.User{}

	if err := c.Bind(&credentials); err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, common.ErrorMessage{
			Message: "invalid payload",
		})
		return
	}
	credentials.Encrypt(salt)

	result := db.First(&user, "username = ? AND password = ?", credentials.Username, credentials.Password)
	if result.RecordNotFound() {
		c.AbortWithStatusJSON(http.StatusUnauthorized, common.ErrorMessage{
			Message: "login denied. Check username or password",
		})
		return
	}

	if result.Error != nil {
		fmt.Printf("%v - \n", result.Error.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, common.ErrorMessage{
			Message: "authentication error. Check logs or contact system admin",
		})
		return
	}

	ttl := time.Now().Add(60 * time.Minute)
	cookie := &Jwt{
		User: user.UUID,
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
