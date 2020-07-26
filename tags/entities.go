package tags

import (
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

// Tag to iterate with database
type Tag struct {
	UUID        uuid.UUID  `gorm:"type:uuid;PRIMARY_KEY" json:"uuid" binding:"-"`
	CreatedAt   time.Time  `json:"createdAt" binding:"-"`
	UpdatedAt   time.Time  `json:"updatedAt" binding:"-"`
	DeletedAt   *time.Time `json:"-" binding:"-"`
	Name        string     `json:"name" binding:"required"`
	Description string     `json:"description"  `
	Owner       uuid.UUID  `gorm:"INDEX,not null" json:"-" `
}

// TableName returns tag table name
func (Tag) TableName() string {
	return "public.tags"
}

// BeforeCreate execute commands before creating a Tag
func (Tag) BeforeCreate(scope *gorm.Scope) (err error) {
	err = scope.SetColumn("UUID", uuid.New())
	return
}
