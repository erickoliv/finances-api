package auth

import (
	"net/http"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/ericktm/olivsoft-golang-api/common"
	"github.com/gin-gonic/gin"
)

// Middleware to validate Authorization Headers
// TODO: Use JWT
func Middleware() gin.HandlerFunc {
	envToken := os.Getenv("APP_TOKEN")

	return func(c *gin.Context) {
		cookie, err := c.Cookie(common.AuthCookie)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, common.ErrorMessage{Message: "auth cookie missing"})
			return
		}

		claims := &Jwt{}
		tkn, err := jwt.ParseWithClaims(cookie, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(envToken), nil
		})

		if err != nil {
			println(err.Error())
			c.AbortWithStatusJSON(http.StatusUnauthorized, common.ErrorMessage{Message: err.Error()})
			return
		}

		if !tkn.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, common.ErrorMessage{Message: "invalid token"})
			return
		}

		c.Set(common.LoggedUser, claims.Username)
		c.Next()
	}
}
