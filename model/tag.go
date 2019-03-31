package model

import (
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

// Tag to iterate with database
type Tag struct {
	BaseModel
	Name        string `json:"name"`
	Description string `json:"description"`
	Owner       uuid.UUID `gorm:"INDEX,not null" json:"owner"`
}

func (t *Tag) BeforeCreate(scope *gorm.Scope) (err error) {
	err = scope.SetColumn("UUID", uuid.New())

	if len(t.Name) == 0 {
		err = errors.New("name cannot be empty")
		return
	}

	//if t.Owner.String() == uuid.New().String() {
	//	err = errors.New("owner cannot be empty")
	//	return
	//}

	return
}


// TableName returns tag table name
func (Tag) TableName() string {
	return "public.tags"
}
