package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

// BaseModel is the base struct for tables
type BaseModel struct {
	UUID      uuid.UUID  `gorm:"type:uuid;PRIMARY_KEY" json:"uuid"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `json:"-"`
}

// BeforeCreate execute commands before creating a BaseModel
func (b BaseModel) BeforeCreate(scope *gorm.Scope) (err error) {
	err = scope.SetColumn("UUID", uuid.New())
	return
}

// IsNew checks if the entity is a new record, without a proper uuid
func (b BaseModel) IsNew() bool {
	return b.UUID == uuid.Nil
}
