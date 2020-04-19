package session

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/erickoliv/finances-api/domain"
	"github.com/erickoliv/finances-api/service"
	"github.com/erickoliv/finances-api/test/entities"
	"github.com/stretchr/testify/assert"
)

func TestNewJWTSigner(t *testing.T) {
	simpleKey := []byte("simple string")

	tests := []struct {
		name string
		ttl  time.Duration
		key  []byte
		want service.Signer
	}{
		{
			name: "creates a new jwt signer",
			key:  simpleKey,
			ttl:  time.Hour,
			want: &jwtSigner{
				key:        simpleKey,
				sessionTTL: time.Hour,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewJWTSigner(tt.key, tt.ttl); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewJWTSigner() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_jwtSigner_SignUser(t *testing.T) {

	key := []byte("simple key")
	ttl := time.Hour
	signer := NewJWTSigner(key, ttl)

	tests := []struct {
		name     string
		signer   service.Signer
		ctx      context.Context
		user     *domain.User
		contains string
		err      error
	}{
		{
			name:     "sign a user uuid",
			signer:   signer,
			ctx:      context.TODO(),
			user:     entities.ValidUser(),
			contains: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
		},
		{
			name:   "validates a nil user",
			signer: signer,
			ctx:    context.TODO(),
			err:    errInvalidUser,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.signer.SignUser(tt.ctx, tt.user)
			assert.Contains(t, got, tt.contains)
			assert.Equal(t, tt.err, err)
		})
	}
}
