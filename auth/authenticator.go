package auth

import (
	"context"
)

type Authenticator interface {
	Login(ctx context.Context, username string, password string) (*User, error)
	Register(ctx context.Context, user *User) error
}
