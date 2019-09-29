package model

import (
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
