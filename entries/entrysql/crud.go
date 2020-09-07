package entrysql

import (
	"context"

	"github.com/erickoliv/finances-api/entries"
	"github.com/erickoliv/finances-api/pkg/http/rest"
	"github.com/erickoliv/finances-api/pkg/querybuilder"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type Repository struct {
	conn *gorm.DB
}

func BuildRepository(conn *gorm.DB) *Repository {
	return &Repository{
		conn,
	}
}

func (repo *Repository) Get(ctx context.Context, pk uuid.UUID, owner uuid.UUID) (*entries.Entry, error) {
	entry := &entries.Entry{}
	query := repo.conn.First(entry, "uuid = ? AND owner = ?", pk, owner)

	if query.Error == nil && entry.UUID == uuid.Nil {
		return nil, errors.New("entry not found")
	}

	return entry, query.Error
}

func (repo *Repository) Query(ctx context.Context, filters *rest.Query) ([]entries.Entry, error) {
	query, err := querybuilder.Build(repo.conn, filters)
	if err != nil {
		return nil, errors.Wrap(err, "entry repository query")
	}

	results := []entries.Entry{}
	query = query.Find(&results)

	return results, query.Error
}

func (repo *Repository) Save(ctx context.Context, entry *entries.Entry) error {
	if entry == nil {
		return errors.New("invalid entry")
	}
	if entry.Owner == uuid.Nil {
		return errors.New("invalid entry. empty owner")
	}

	return repo.conn.Save(entry).Error
}

func (repo *Repository) Delete(ctx context.Context, id uuid.UUID, user uuid.UUID) error {
	query := repo.conn.Where("uuid = ? AND owner = ?", id, user).Delete(&entries.Entry{})
	return query.Error
}
