package model

import (
	"time"

	"github.com/google/uuid"
)

// Entry to iterate with database
type Entry struct {
	BaseModel
	Date        time.Time `json:"date" binding:"required"`
	Type        bool      `json:"type" binding:"required"`
	Pending     bool      `json:"pending"`
	Name        string    `json:"name" binding:"required"`
	Value       float64   `json:"value" binding:"required"`
	Description string    `json:"description"  `
	IsTransfer  bool      `json:"isTransfer"`
	Origin      uuid.UUID `json:"origin"`
	Category    uuid.UUID `json:"category" binding:"required" gorm:"INDEX,not null"`
	Account     uuid.UUID `json:"account" binding:"required" gorm:"INDEX,not null"`
	Owner       uuid.UUID `gorm:"INDEX,not null" json:"owner" `
}

// TableName returns Entry table name
func (Entry) TableName() string {
	return "public.entries"
}
