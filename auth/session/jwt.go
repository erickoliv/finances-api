package session

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func NewJWTSigner(key []byte, ttl time.Duration) *jwtSigner {
	return &jwtSigner{
		key:        key,
		sessionTTL: ttl,
	}
}

func (signer *jwtSigner) SignUser(identifier string) (string, error) {
	if len(identifier) == 0 {
		return "", errEmptyIdentifier
	}
	if len(signer.key) == 0 {
		return "", errInvalidKey
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &claim{
		User: identifier,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(signer.sessionTTL).Unix(),
		},
	})

	return token.SignedString(signer.key)
}

func (signer *jwtSigner) Validate(token string) (string, error) {
	claims := &claim{}
	_, err := jwt.ParseWithClaims(token, claims, signer.keyFunc)
	if err != nil {
		return "", fmt.Errorf("jwt validate error: %s", err.Error())
	}
	return claims.User, nil
}

func (signer *jwtSigner) keyFunc(token *jwt.Token) (interface{}, error) {
	return []byte(signer.key), nil
}
