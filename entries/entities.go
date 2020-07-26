package entries

import (
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

// Entry to iterate with database
type Entry struct {
	UUID        uuid.UUID  `gorm:"type:uuid;PRIMARY_KEY" json:"uuid" binding:"-"`
	CreatedAt   time.Time  `json:"createdAt" binding:"-"`
	UpdatedAt   time.Time  `json:"updatedAt" binding:"-"`
	DeletedAt   *time.Time `json:"-" binding:"-"`
	Date        time.Time  `json:"date" binding:"required"`
	Type        bool       `json:"type" binding:"required"`
	Pending     bool       `json:"pending"`
	Name        string     `json:"name" binding:"required"`
	Value       float64    `json:"value" binding:"required"`
	Description string     `json:"description"  `
	IsTransfer  bool       `json:"isTransfer"`
	Origin      uuid.UUID  `json:"origin"`
	Category    uuid.UUID  `json:"category" binding:"required" gorm:"INDEX,not null"`
	Account     uuid.UUID  `json:"account" binding:"required" gorm:"INDEX,not null"`
	Owner       uuid.UUID  `gorm:"INDEX,not null" json:"owner" `
}

// TableName returns Entry table name
func (Entry) TableName() string {
	return "public.entries"
}

// BeforeCreate execute commands before creating a Entry
func (Entry) BeforeCreate(scope *gorm.Scope) error {
	return scope.SetColumn("UUID", uuid.New())
}

// EntryTag is a association between Entry and Tag entities
type EntryTag struct {
	UUID      uuid.UUID  `gorm:"type:uuid;PRIMARY_KEY" json:"uuid" binding:"-"`
	CreatedAt time.Time  `json:"createdAt" binding:"-"`
	UpdatedAt time.Time  `json:"updatedAt" binding:"-"`
	DeletedAt *time.Time `json:"-" binding:"-"`
	Entry     uuid.UUID  `gorm:"INDEX,not null" json:"entry" `
	Tag       uuid.UUID  `gorm:"INDEX,not null" json:"tag" `
}

// TableName returns Entry table name
func (EntryTag) TableName() string {
	return "public.entry_tags"
}

// BeforeCreate execute command before creating a EntryTag
func (EntryTag) BeforeCreate(scope *gorm.Scope) (err error) {
	err = scope.SetColumn("UUID", uuid.New())
	return
}
