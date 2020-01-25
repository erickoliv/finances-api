package auth

import (
	"net/http"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/erickoliv/finances-api/domain"
	"github.com/gin-gonic/gin"
)

// Middleware to validate Authorization Headers
// TODO: Use JWT
func Middleware() gin.HandlerFunc {
	envToken := os.Getenv("APP_TOKEN")

	return func(c *gin.Context) {
		cookie, err := c.Cookie(domain.AuthCookie)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, domain.ErrorMessage{Message: "auth cookie missing"})
			return
		}

		claims := &Jwt{}
		tkn, err := jwt.ParseWithClaims(cookie, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(envToken), nil
		})

		if err != nil {
			println(err.Error())
			c.AbortWithStatusJSON(http.StatusUnauthorized, domain.ErrorMessage{Message: err.Error()})
			return
		}

		if !tkn.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, domain.ErrorMessage{Message: "invalid token"})
			return
		}

		c.Set(domain.LoggedUser, claims.User)
		c.Next()
	}
}
