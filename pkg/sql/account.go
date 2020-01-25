package sql

import (
	"context"
	"github.com/erickoliv/finances-api/domain"
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

func (repo accountRepo) Get(ctx context.Context, pk uuid.UUID, owner uuid.UUID) (*domain.Account, error) {
	result := &domain.Account{}

	status := repo.conn.Where("uuid = ?", pk).Where("owner = ?", owner)
	status.First(result)

	return result, status.Error
}

func (repo accountRepo) Filter(ctx context.Context, query domain.Query) ([]domain.Account, error) {
	result := []domain.Account{}
	status := BuildQuery(repo.conn, query).Find(&result)

	return result, status.Error
}

func (repo accountRepo) Save(ctx context.Context, account *domain.Account) error {
	return repo.conn.Save(account).Error
}

func (repo accountRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return ctx.Err()
}
