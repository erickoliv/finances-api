package model

import (
	"github.com/google/uuid"
)

// Category to iterate with database
type Category struct {
	BaseModel
	Name        string    `json:"name" binding:"required"`
	Description string    `json:"description"  `
	Parent      uuid.UUID `json:"parent" gorm:"INDEX"`
	Owner       uuid.UUID `gorm:"INDEX,not null" json:"owner" `
}

// TableName returns Category table name
func (Category) TableName() string {
	return "public.categories"
}
