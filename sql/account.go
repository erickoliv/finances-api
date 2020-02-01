package sql

import (
	"context"

	"github.com/erickoliv/finances-api/domain"
	"github.com/erickoliv/finances-api/pkg/http/rest"
	"github.com/erickoliv/finances-api/repository"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

type accountRepo struct {
	conn *gorm.DB
}

func MakeAccounts(conn *gorm.DB) repository.AccountService {
	return accountRepo{
		conn,
	}
}

func (repo accountRepo) Get(ctx context.Context, pk uuid.UUID, owner uuid.UUID) (*domain.Account, error) {
	account := &domain.Account{}
	query := repo.conn.First(account, "uuid = ? AND owner = ?", pk, owner)

	return account, query.Error
}

func (repo accountRepo) Query(ctx context.Context, filters *rest.Query) ([]*domain.Account, error) {
	query, err := BuildQuery(repo.conn, filters)
	if err != nil {
		return nil, errors.Wrap(err, "account repository query")
	}

	results := []*domain.Account{}
	query = query.Find(&results)

	return results, query.Error
}

func (repo accountRepo) Save(ctx context.Context, account *domain.Account) error {
	if account == nil {
		return errors.New("invalid account")
	}

	return repo.conn.Save(account).Error
}

func (repo accountRepo) Delete(ctx context.Context, id uuid.UUID, user uuid.UUID) error {
	query := repo.conn.Where("uuid = ? AND owner = ?", id, user).Delete(&domain.Account{})
	return query.Error
}
