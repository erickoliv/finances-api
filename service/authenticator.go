package service

import (
	"context"

	"github.com/erickoliv/finances-api/domain"
)

// Authenticator defines the methods to handle application authentication
type Authenticator interface {
	Login(ctx context.Context, username string, password string) (*domain.User, error)
	Register(ctx context.Context, user *domain.User) error
}
