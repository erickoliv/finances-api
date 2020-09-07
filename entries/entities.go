package entries

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Entry to iterate with database
type Entry struct {
	UUID        uuid.UUID      `gorm:"type:uuid;primaryKey" json:"uuid" binding:"-"`
	CreatedAt   time.Time      `json:"createdAt" binding:"-"`
	UpdatedAt   time.Time      `json:"updatedAt" binding:"-"`
	DeletedAt   gorm.DeletedAt `json:"-" binding:"-"`
	Date        time.Time      `json:"date" binding:"required"`
	Type        bool           `json:"type" binding:"required"`
	Pending     bool           `json:"pending"`
	Name        string         `json:"name" binding:"required"`
	Value       float64        `json:"value" binding:"required"`
	Description string         `json:"description"  `
	IsTransfer  bool           `json:"isTransfer"`
	Origin      uuid.UUID      `json:"origin"`
	Category    uuid.UUID      `json:"category" binding:"required" gorm:"index,not null"`
	Account     uuid.UUID      `json:"account" binding:"required" gorm:"index,not null"`
	Owner       uuid.UUID      `gorm:"index,not null" json:"owner" `
}

// TableName returns Entry table name
func (Entry) TableName() string {
	return "public.entries"
}

// BeforeCreate execute commands before creating a Entry
func (e *Entry) BeforeCreate(scope *gorm.DB) error {
	e.UUID = uuid.New()

	return nil
}

// EntryTag is a association between Entry and Tag entities
type EntryTag struct {
	UUID      uuid.UUID  `gorm:"type:uuid;primaryKey" json:"uuid" binding:"-"`
	CreatedAt time.Time  `json:"createdAt" binding:"-"`
	UpdatedAt time.Time  `json:"updatedAt" binding:"-"`
	DeletedAt *time.Time `json:"-" binding:"-"`
	Entry     uuid.UUID  `gorm:"index:entry_tag_entry;not null" json:"entry" `
	Tag       uuid.UUID  `gorm:"index:entry_tag_tag;not null" json:"tag" `
}

// TableName returns Entry table name
func (EntryTag) TableName() string {
	return "public.entry_tags"
}

// BeforeCreate execute command before creating a EntryTag
func (e *EntryTag) BeforeCreate(scope *gorm.DB) (err error) {
	e.UUID = uuid.New()

	return
}
