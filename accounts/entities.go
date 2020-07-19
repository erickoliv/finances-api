package accounts

import (
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

// Account to iterate with database
type Account struct {
	UUID        uuid.UUID  `gorm:"type:uuid;PRIMARY_KEY" json:"uuid" binding:"-"`
	CreatedAt   time.Time  `json:"createdAt" binding:"-"`
	UpdatedAt   time.Time  `json:"updatedAt" binding:"-"`
	DeletedAt   *time.Time `json:"-" binding:"-"`
	Name        string     `json:"name" binding:"required"`
	Description string     `json:"description"  `
	Owner       uuid.UUID  `gorm:"INDEX,not null" json:"-" `
}

// TableName returns Account table name
func (Account) TableName() string {
	return "public.accounts"
}

// BeforeCreate execute commands before creating a Account
func (Account) BeforeCreate(scope *gorm.Scope) (err error) {
	err = scope.SetColumn("UUID", uuid.New())
	return
}
