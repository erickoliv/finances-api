package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_encrypt(t *testing.T) {
	tests := []struct {
		name string
		pass string
		salt string
		want string
	}{
		{
			name: "encrypts a password with a salt and returns a string",
			pass: "123456",
			salt: "958d51602bbfbd18b2a084ba848a827c29952bfef170c936419b0922994c0589",
			want: "655b7974d95ca0e9a4c4b84444eecb9f61920064f8e542e626ce576e1bf11d2f",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, encrypt(tt.pass, tt.salt))
		})
	}
}
