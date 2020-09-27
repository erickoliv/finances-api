package jwtsession

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewJWTSigner(t *testing.T) {
	simpleKey := []byte("simple string")

	tests := []struct {
		name string
		ttl  time.Duration
		key  []byte
		want *JWTSigner
	}{
		{
			name: "creates a new jwt signer",
			key:  simpleKey,
			ttl:  time.Hour,
			want: &JWTSigner{
				key:        simpleKey,
				sessionTTL: time.Hour,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			signer := NewJWTSigner(tt.key, tt.ttl)
			assert.Equal(t, signer, tt.want)
		})
	}
}

func Test_SignAndValidateToken(t *testing.T) {
	signKey := []byte(time.Now().String())

	tests := []struct {
		name         string
		key          []byte
		sessionTTL   time.Duration
		identifier   string
		wantedSigner *JWTSigner
		signErr      error
		validateErr  error
	}{
		{
			name:       "successfully signs and then validate a user token",
			key:        signKey,
			sessionTTL: time.Hour,
			identifier: uuid.New().String(),
			wantedSigner: &JWTSigner{
				key:        signKey,
				sessionTTL: time.Hour,
			},
		},
		{
			name:       "successfully signs a user identifier but fails to validate token due to expiration",
			key:        signKey,
			sessionTTL: time.Millisecond,
			identifier: uuid.New().String(),
			wantedSigner: &JWTSigner{
				key:        signKey,
				sessionTTL: time.Millisecond,
			},
			validateErr: errors.New("jwt validate error: token is expired by 1s"),
		},
		{
			name:       "error to sign due to empty user identifier",
			key:        signKey,
			sessionTTL: time.Hour,
			identifier: "",
			wantedSigner: &JWTSigner{
				key:        signKey,
				sessionTTL: time.Hour,
			},
			signErr: errEmptyIdentifier,
		},
		{
			name:       "error to sign due to empty key",
			sessionTTL: time.Hour,
			identifier: uuid.New().String(),
			wantedSigner: &JWTSigner{
				sessionTTL: time.Hour,
			},
			signErr: errInvalidKey,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			signer := NewJWTSigner(tt.key, tt.sessionTTL)
			assert.Equal(t, tt.wantedSigner, signer)

			token, err := signer.SignUser(tt.identifier)
			assert.Equal(t, tt.signErr, err)
			if err != nil {
				return
			}

			time.Sleep(time.Second)
			identifier, err := signer.Validate(token)
			assert.Equal(t, tt.validateErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tt.identifier, identifier)
		})
	}
}
