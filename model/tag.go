package model

import (
	"github.com/google/uuid"
)

// Tag to iterate with database
type Tag struct {
	BaseModel
	Name        string    `json:"name" binding:"required"`
	Description string    `json:"description"  `
	Owner       uuid.UUID `gorm:"INDEX,not null" json:"owner" `
}

// TableName returns tag table name
func (Tag) TableName() string {
	return "public.tags"
}
