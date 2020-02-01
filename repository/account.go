package repository

import (
	"context"

	"github.com/erickoliv/finances-api/domain"
	"github.com/erickoliv/finances-api/pkg/http/rest"
	"github.com/google/uuid"
)

type AccountService interface {
	Delete(ctx context.Context, pk uuid.UUID, owner uuid.UUID) error
	Get(ctx context.Context, pk uuid.UUID, owner uuid.UUID) (*domain.Account, error)
	Save(ctx context.Context, account *domain.Account) error
	Query(ctx context.Context, filters *rest.Query) ([]*domain.Account, error)
}
