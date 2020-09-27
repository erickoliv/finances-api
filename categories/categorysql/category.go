package categorysql

import (
	"context"

	"github.com/erickoliv/finances-api/categories"
	"github.com/erickoliv/finances-api/pkg/http/rest"
	"github.com/erickoliv/finances-api/pkg/querybuilder"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type CategoryRepo struct {
	conn *gorm.DB
}

func BuildRepository(conn *gorm.DB) *CategoryRepo {
	return &CategoryRepo{
		conn,
	}
}

func (repo *CategoryRepo) Get(ctx context.Context, pk uuid.UUID, owner uuid.UUID) (*categories.Category, error) {
	category := &categories.Category{}
	query := repo.conn.First(category, "uuid = ? AND owner = ?", pk, owner)

	if query.Error == nil && category.UUID == uuid.Nil {
		return nil, errors.New("category not found")
	}

	return category, query.Error
}

func (repo *CategoryRepo) Query(ctx context.Context, filters *rest.Query) ([]categories.Category, error) {
	query, err := querybuilder.Build(repo.conn, filters)
	if err != nil {
		return nil, errors.Wrap(err, "category repository query")
	}

	results := []categories.Category{}
	query = query.Find(&results)

	return results, query.Error
}

func (repo *CategoryRepo) Save(ctx context.Context, category *categories.Category) error {
	if category == nil {
		return errors.New("invalid category")
	}
	if category.Owner == uuid.Nil {
		return errors.New("invalid category. empty owner")
	}

	return repo.conn.Save(category).Error
}

func (repo *CategoryRepo) Delete(ctx context.Context, id uuid.UUID, user uuid.UUID) error {
	query := repo.conn.Where("uuid = ? AND owner = ?", id, user).Delete(&categories.Category{})
	return query.Error
}
