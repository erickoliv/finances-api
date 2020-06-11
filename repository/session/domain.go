package session

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type claim struct {
	User string
	jwt.StandardClaims
}

type jwtSigner struct {
	key        []byte
	sessionTTL time.Duration
}

var (
	errInvalidKey      = errors.New("invalid empty key")
	errInvalidToken    = errors.New("token is invalid or expired")
	errEmptyIdentifier = errors.New("empty user identifier")
)
