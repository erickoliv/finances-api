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

// Jwt struct to generate and validate jtw tokens
type Jwt struct {
	User uuid.UUID `json:"user"`
	jwt.StandardClaims
}

func encrypt(pass string, salt string) string {
	hash := sha256.Sum256([]byte(pass + salt))
	return fmt.Sprintf("%x", hash)
}
