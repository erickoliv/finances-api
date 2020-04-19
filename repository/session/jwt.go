package session

import (
	"context"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/erickoliv/finances-api/domain"
	"github.com/erickoliv/finances-api/service"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type claim struct {
	User uuid.UUID
	jwt.StandardClaims
}

type jwtSigner struct {
	key        []byte
	sessionTTL time.Duration
}

var (
	errInvalidKey   = errors.New("invalid key")
	errInvalidToken = errors.New("token is invalid or expired")
	errInvalidUser  = errors.New("invalid user")
)

func NewJWTSigner(key []byte, ttl time.Duration) service.Signer {
	return &jwtSigner{
		key:        key,
		sessionTTL: ttl,
	}
}

func (signer *jwtSigner) SignUser(ctx context.Context, user *domain.User) (string, error) {
	if user == nil {
		return "", errInvalidUser
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &claim{
		User: user.UUID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(signer.sessionTTL).Unix(),
		},
	})

	return token.SignedString(signer.key)
}

func (signer *jwtSigner) Validate(ctx context.Context, token string) (uuid.UUID, error) {
	claims := &claim{}

	tkn, err := jwt.ParseWithClaims(token, claims, signer.keyFunc)
	if err != nil {
		return uuid.Nil, err
	}

	if !tkn.Valid {
		return uuid.Nil, errInvalidToken
	}

	return claims.User, nil
}

func (signer *jwtSigner) keyFunc(token *jwt.Token) (interface{}, error) {
	if len(signer.key) == 0 {
		return nil, errInvalidKey
	}
	return []byte(signer.key), nil
}
