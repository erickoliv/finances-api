package entries

import (
	"context"

	"github.com/erickoliv/finances-api/pkg/http/rest"
	"github.com/google/uuid"
)

// Repository interface defines methods to manipute user entries
type Repository interface {
	Delete(ctx context.Context, pk uuid.UUID, owner uuid.UUID) error
	Get(ctx context.Context, pk uuid.UUID, owner uuid.UUID) (*Entry, error)
	Save(ctx context.Context, row *Entry) error
	Query(ctx context.Context, filters *rest.Query) ([]Entry, error)
}
