package categories

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Category to iterate with database
type Category struct {
	UUID        uuid.UUID      `gorm:"type:uuid;primaryKey" json:"uuid" binding:"-"`
	CreatedAt   time.Time      `json:"createdAt" binding:"-"`
	UpdatedAt   time.Time      `json:"updatedAt" binding:"-"`
	DeletedAt   gorm.DeletedAt `json:"-" binding:"-"`
	Name        string         `json:"name" binding:"required"`
	Description string         `json:"description"  `
	Parent      uuid.UUID      `json:"parent" gorm:"index:category_parent"`
	Owner       uuid.UUID      `gorm:"index:category_owner;not null" json:"owner" `
}

// TableName returns Category table name
func (Category) TableName() string {
	return "public.categories"
}

// BeforeCreate execute commands before creating a Category
func (c *Category) BeforeCreate(scope *gorm.DB) (err error) {
	c.UUID = uuid.New()

	return
}
