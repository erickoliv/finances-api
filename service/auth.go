package service

import (
	"context"

	"github.com/erickoliv/finances-api/domain"
	"github.com/google/uuid"
)

// Authenticator defines the methods to handle application authentication
type Authenticator interface {
	Login(context.Context, string, string) (*domain.User, error)
	Register(context.Context, *domain.User) error
}

// Signer defines the methods to sign authentication data
type Signer interface {
	SignUser(context.Context, *domain.User) (string, error)
	Validate(context.Context, string) (uuid.UUID, error)
}
