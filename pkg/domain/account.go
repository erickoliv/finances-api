package domain

import (
	"context"
	"github.com/google/uuid"
)

// Account to iterate with database
type Account struct {
	BaseModel
	Name        string    `json:"name" binding:"required"`
	Description string    `json:"description"  `
	Owner       uuid.UUID `gorm:"INDEX,not null" json:"owner" `
}

// TableName returns Account table name
func (Account) TableName() string {
	return "public.accounts"
}

type AccountRepository interface {
	Delete(context.Context, uuid.UUID) error
	Filter(context.Context, QueryData) ([]Account, error)
	Get(context.Context, uuid.UUID) (Account, error)
	Save(context.Context, Account) error
}
