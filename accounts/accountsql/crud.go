package accountsql

import (
	"context"

	"github.com/erickoliv/finances-api/accounts"
	"github.com/erickoliv/finances-api/pkg/http/rest"
	"github.com/erickoliv/finances-api/pkg/querybuilder"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

type accountRepo struct {
	conn *gorm.DB
}

func MakeAccounts(conn *gorm.DB) *accountRepo {
	return &accountRepo{
		conn,
	}
}

func (repo *accountRepo) Get(ctx context.Context, pk uuid.UUID, owner uuid.UUID) (*accounts.Account, error) {
	account := &accounts.Account{}
	query := repo.conn.First(account, "uuid = ? AND owner = ?", pk, owner)

	if query.Error == nil && account.UUID == uuid.Nil {
		return nil, nil
	}

	return account, query.Error
}

func (repo *accountRepo) Query(ctx context.Context, filters *rest.Query) ([]*accounts.Account, error) {
	query, err := querybuilder.Build(repo.conn, filters)
	if err != nil {
		return nil, errors.Wrap(err, "account repository query")
	}

	results := []*accounts.Account{}
	query = query.Find(&results)

	return results, query.Error
}

func (repo *accountRepo) Save(ctx context.Context, account *accounts.Account) error {
	if account == nil {
		return errors.New("invalid account")
	}

	return repo.conn.Save(account).Error
}

func (repo *accountRepo) Delete(ctx context.Context, id uuid.UUID, user uuid.UUID) error {
	query := repo.conn.Where("uuid = ? AND owner = ?", id, user).Delete(&accounts.Account{})
	return query.Error
}
