package model

import (
	"time"

	"github.com/google/uuid"
)

// BaseModel is the base struct for tables
type BaseModel struct {
	UUID      uuid.UUID `gorm:"type:uuid;PRIMARY_KEY" json:"uuid"`
	CreatedAt time.Time `json:"created-at"`
	UpdatedAt time.Time `json:"updated-at"`
	DeletedAt *time.Time `json:"deleted-at"`
}
