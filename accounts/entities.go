package accounts

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Account to iterate with database
type Account struct {
	UUID        uuid.UUID      `gorm:"type:uuid;primaryKey" json:"uuid" binding:"-"`
	CreatedAt   time.Time      `json:"createdAt" binding:"-"`
	UpdatedAt   time.Time      `json:"updatedAt" binding:"-"`
	DeletedAt   gorm.DeletedAt `json:"-" binding:"-"`
	Name        string         `json:"name" binding:"required"`
	Description string         `json:"description"  `
	Owner       uuid.UUID      `gorm:"index,not null" json:"-" `
}

// TableName returns Account table name
func (Account) TableName() string {
	return "public.accounts"
}

// BeforeCreate execute commands before creating a Account
func (a *Account) BeforeCreate(scope *gorm.DB) (err error) {
	a.UUID = uuid.New()

	return
}
