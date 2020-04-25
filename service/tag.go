package service

import (
	"context"

	"github.com/erickoliv/finances-api/domain"
	"github.com/erickoliv/finances-api/pkg/http/rest"
	"github.com/google/uuid"
)

// Tag interface defines methods to manipute user entry's tags
type Tag interface {
	Delete(ctx context.Context, pk uuid.UUID, owner uuid.UUID) error
	Get(ctx context.Context, pk uuid.UUID, owner uuid.UUID) (*domain.Tag, error)
	Save(ctx context.Context, tag *domain.Tag) error
	Query(ctx context.Context, filters *rest.Query) ([]*domain.Tag, error)
}
