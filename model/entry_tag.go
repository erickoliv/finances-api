package model

import (
	"github.com/google/uuid"
)

// EntryTag is a association between Entry and Tag entities
type EntryTag struct {
	BaseModel
	Entry uuid.UUID `gorm:"INDEX,not null" json:"entry" `
	Tag   uuid.UUID `gorm:"INDEX,not null" json:"tag" `
}

// TableName returns Entry table name
func (EntryTag) TableName() string {
	return "public.entry_tags"
}
