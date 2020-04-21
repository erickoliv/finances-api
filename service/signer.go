package service

import (
	"context"

	"github.com/erickoliv/finances-api/domain"
	"github.com/google/uuid"
)

// Signer defines the methods to sign authentication data
type Signer interface {
	SignUser(ctx context.Context, user *domain.User) (string, error)
	Validate(ctx context.Context, token string) (uuid.UUID, error)
}
