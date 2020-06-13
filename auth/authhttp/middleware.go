package authhttp

import (
	"net/http"

	"github.com/erickoliv/finances-api/auth"
	"github.com/erickoliv/finances-api/domain"
	"github.com/gin-gonic/gin"
)

// Middleware to validate JWT authentication cookie
func Middleware(signer auth.SessionSigner) gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie(auth.AuthCookie)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "auth cookie missing"})
			return
		}

		user, err := signer.Validate(cookie)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
			return
		}

		c.Set(domain.LoggedUser, user)
		c.Next()
	}
}
