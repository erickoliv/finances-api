package auth

import (
	"crypto/sha256"
	"fmt"

	"github.com/dgrijalva/jwt-go"
)

// Credentials struct used to parse login requests
type Credentials struct {
	Username string `json:"username" binding:"required" `
	Password string `json:"password" binding:"required" `
}

// Encrypt password using sha256
func (c *Credentials) Encrypt() {
	hash := sha256.Sum256([]byte(c.Password + "TODO: add salt"))
	c.Password = fmt.Sprintf("%x", hash)
}

// JWT
type Jwt struct {
	Username string `json:"username"`
	jwt.StandardClaims
}
