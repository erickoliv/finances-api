package tags

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Tag to iterate with database
type Tag struct {
	UUID        uuid.UUID      `gorm:"type:uuid;primaryKey" json:"uuid" binding:"-"`
	CreatedAt   time.Time      `json:"createdAt" binding:"-"`
	UpdatedAt   time.Time      `json:"updatedAt" binding:"-"`
	DeletedAt   gorm.DeletedAt `json:"-" binding:"-"`
	Name        string         `json:"name" binding:"required"`
	Description string         `json:"description"  `
	Owner       uuid.UUID      `gorm:"index:tag_owner;not null" json:"-" `
}

// TableName returns tag table name
func (Tag) TableName() string {
	return "public.tags"
}

// BeforeCreate execute commands before creating a Tag
func (t *Tag) BeforeCreate(scope *gorm.DB) (err error) {
	t.UUID = uuid.New()
	return
}
