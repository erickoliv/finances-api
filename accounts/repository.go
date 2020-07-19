package accounts

import (
	"context"

	"github.com/erickoliv/finances-api/pkg/http/rest"
	"github.com/google/uuid"
)

// Repository interface defines methods to manipute user's accounts
type Repository interface {
	Delete(ctx context.Context, pk uuid.UUID, owner uuid.UUID) error
	Get(ctx context.Context, pk uuid.UUID, owner uuid.UUID) (*Account, error)
	Save(ctx context.Context, account *Account) error
	Query(ctx context.Context, filters *rest.Query) ([]*Account, error)
}
