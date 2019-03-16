package model

import (
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

// Tag to iterate with database
type Tag struct {
	BaseModel
	Name        string `json:"name"`
	Description string `json:"description"`
	Owner       uuid.UUID `gorm:"INDEX,not null" json:"owner"`
}

func (t *Tag) BeforeCreate(scope *gorm.Scope) error {
	return scope.SetColumn("UUID", uuid.New())
}

// TableName returns tag table name
func (Tag) TableName() string {
	return "public.tags"
}
