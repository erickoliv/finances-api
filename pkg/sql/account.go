package sql

import (
	"context"
	"github.com/ericktm/olivsoft-golang-api/pkg/domain"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

type accountRepo struct {
	conn *gorm.DB
}

func MakeAccounts(conn *gorm.DB) domain.AccountRepository {
	return accountRepo{
		conn,
	}
}

func (repo accountRepo) Get(ctx context.Context, id uuid.UUID) (domain.Account, error) {
	result := domain.Account{}

	return result, ctx.Err()
}

func (repo accountRepo) Filter(ctx context.Context, query domain.QueryData) ([]domain.Account, error) {
	result := []domain.Account{}
	status := BuildQuery(repo.conn, query).Find(&result)

	return result, status.Error
}

func (repo accountRepo) Save(ctx context.Context, account domain.Account) error {
	return ctx.Err()
}

func (repo accountRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return ctx.Err()
}
