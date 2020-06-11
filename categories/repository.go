package categories

import (
	"context"

	"github.com/erickoliv/finances-api/pkg/http/rest"
	"github.com/google/uuid"
)

// Repository interface defines methods to manipute user categories
type Repository interface {
	Delete(ctx context.Context, pk uuid.UUID, owner uuid.UUID) error
	Get(ctx context.Context, pk uuid.UUID, owner uuid.UUID) (*Category, error)
	Save(ctx context.Context, row *Category) error
	Query(ctx context.Context, filters *rest.Query) ([]Category, error)
}
