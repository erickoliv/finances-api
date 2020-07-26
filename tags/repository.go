package tags

import (
	"context"

	"github.com/erickoliv/finances-api/pkg/http/rest"
	"github.com/google/uuid"
)

// Repository defines methods to manipute user entry's tags
type Repository interface {
	Delete(ctx context.Context, pk uuid.UUID, owner uuid.UUID) error
	Get(ctx context.Context, pk uuid.UUID, owner uuid.UUID) (*Tag, error)
	Save(ctx context.Context, tag *Tag) error
	Query(ctx context.Context, filters *rest.Query) ([]Tag, error)
}
