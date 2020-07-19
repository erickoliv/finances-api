package auth

import (
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

// AuthCookie is cookie name used in autentication api, maybe use a injected parameter instead
const AuthCookie string = "olivsoftauth"

// User domain/database representation
type User struct {
	UUID      uuid.UUID  `gorm:"type:uuid;PRIMARY_KEY" json:"uuid" binding:"-"`
	CreatedAt time.Time  `json:"createdAt" binding:"-"`
	UpdatedAt time.Time  `json:"updatedAt" binding:"-"`
	DeletedAt *time.Time `json:"-" binding:"-"`
	FirstName string     `json:"firstName" binding:"required"`
	LastName  string     `json:"lastName"`
	Email     string     `json:"email" gorm:"UNIQUE" binding:"required"`
	Username  string     `json:"username" gorm:"UNIQUE" binding:"required"`
	Password  string     `json:"password"`
	Active    bool       `json:"active"  `
}

// TableName set user table name
func (User) TableName() string {
	return "public.users"
}

// BeforeCreate execute commands before creating User
func (User) BeforeCreate(scope *gorm.Scope) (err error) {
	err = scope.SetColumn("UUID", uuid.New())
	return
}
