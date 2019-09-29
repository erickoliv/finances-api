package auth

import (
	"crypto/sha256"
	"fmt"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

// Credentials struct used to parse login requests
type Credentials struct {
	Username string `json:"username" binding:"required" `
	Password string `json:"password" binding:"required" `
}

// Encrypt password using sha256
func (c *Credentials) Encrypt(salt string) {
	hash := sha256.Sum256([]byte(c.Password + salt))
	c.Password = fmt.Sprintf("%x", hash)
}

// Jwt struct to generate and validate jtw tokens
type Jwt struct {
	User uuid.UUID `json:"user"`
	jwt.StandardClaims
}
