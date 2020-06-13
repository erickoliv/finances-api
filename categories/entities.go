package categories

import (
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

// Category to iterate with database
type Category struct {
	UUID        uuid.UUID  `gorm:"type:uuid;PRIMARY_KEY" json:"uuid" binding:"-"`
	CreatedAt   time.Time  `json:"createdAt" binding:"-"`
	UpdatedAt   time.Time  `json:"updatedAt" binding:"-"`
	DeletedAt   *time.Time `json:"-" binding:"-"`
	Name        string     `json:"name" binding:"required"`
	Description string     `json:"description"  `
	Parent      uuid.UUID  `json:"parent" gorm:"INDEX"`
	Owner       uuid.UUID  `gorm:"INDEX,not null" json:"owner" `
}

// TableName returns Category table name
func (Category) TableName() string {
	return "public.categories"
}

// BeforeCreate execute commands before creating a Category
func (c Category) BeforeCreate(scope *gorm.Scope) (err error) {
	err = scope.SetColumn("UUID", uuid.New())
	return
}
