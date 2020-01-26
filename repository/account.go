package repository

import (
	"context"

	"github.com/erickoliv/finances-api/domain"
	"github.com/erickoliv/finances-api/pkg/http/rest"
	"github.com/google/uuid"
)

type AccountService interface {
	Delete(context.Context, uuid.UUID) error
	Filter(context.Context, rest.Query) ([]domain.Account, error)
	Get(context.Context, uuid.UUID, uuid.UUID) (*domain.Account, error)
	Save(context.Context, *domain.Account) error
}
