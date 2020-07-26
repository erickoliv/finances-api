package tagsql

import (
	"context"

	"github.com/erickoliv/finances-api/pkg/http/rest"
	"github.com/erickoliv/finances-api/pkg/querybuilder"
	"github.com/erickoliv/finances-api/tags"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

type tagRepo struct {
	conn *gorm.DB
}

func BuildTagRepository(conn *gorm.DB) *tagRepo {
	return &tagRepo{
		conn,
	}
}

func (repo *tagRepo) Get(ctx context.Context, pk uuid.UUID, owner uuid.UUID) (*tags.Tag, error) {
	tag := &tags.Tag{}
	query := repo.conn.First(tag, "uuid = ? AND owner = ?", pk, owner)

	if query.Error == nil && tag.UUID == uuid.Nil {
		return nil, nil
	}

	return tag, query.Error
}

func (repo *tagRepo) Query(ctx context.Context, filters *rest.Query) ([]tags.Tag, error) {
	query, err := querybuilder.Build(repo.conn, filters)
	if err != nil {
		return nil, errors.Wrap(err, "tag repository query")
	}

	results := []tags.Tag{}
	query = query.Find(&results)

	return results, query.Error
}

func (repo tagRepo) Save(ctx context.Context, tag *tags.Tag) error {
	if tag == nil {
		return errors.New("invalid tag")
	}
	if tag.Owner == uuid.Nil {
		return errors.New("invalid tag. empty owner")
	}

	return repo.conn.Save(tag).Error
}

func (repo tagRepo) Delete(ctx context.Context, id uuid.UUID, user uuid.UUID) error {
	query := repo.conn.Where("uuid = ? AND owner = ?", id, user).Delete(&tags.Tag{})
	return query.Error
}
