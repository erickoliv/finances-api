package sql

import (
	"context"

	"github.com/erickoliv/finances-api/domain"
	"github.com/erickoliv/finances-api/pkg/http/rest"
	"github.com/erickoliv/finances-api/service"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

type tagRepo struct {
	conn *gorm.DB
}

func BuildTagRepository(conn *gorm.DB) service.Tag {
	return tagRepo{
		conn,
	}
}

func (repo tagRepo) Get(ctx context.Context, pk uuid.UUID, owner uuid.UUID) (*domain.Tag, error) {
	tag := &domain.Tag{}
	query := repo.conn.First(tag, "uuid = ? AND owner = ?", pk, owner)

	if query.Error == nil && tag.IsNew() {
		return nil, nil
	}

	return tag, query.Error
}

func (repo tagRepo) Query(ctx context.Context, filters *rest.Query) ([]*domain.Tag, error) {
	query, err := BuildQuery(repo.conn, filters)
	if err != nil {
		return nil, errors.Wrap(err, "tag repository query")
	}

	results := []*domain.Tag{}
	query = query.Find(&results)

	return results, query.Error
}

func (repo tagRepo) Save(ctx context.Context, tag *domain.Tag) error {
	if tag == nil {
		return errors.New("invalid tag")
	}
	if tag.Owner == uuid.Nil {
		return errors.New("invalid tag. empty owner")
	}

	return repo.conn.Save(tag).Error
}

func (repo tagRepo) Delete(ctx context.Context, id uuid.UUID, user uuid.UUID) error {
	query := repo.conn.Where("uuid = ? AND owner = ?", id, user).Delete(&domain.Tag{})
	return query.Error
}
