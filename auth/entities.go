package auth

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	// AuthCookie is cookie name used in autentication api, maybe use a injected parameter instead
	AuthCookie string = "olivsoftauth"

	// LoggedUser contains the UUID for the current user inside context
	LoggedUser string = "current-logged-user"
)

// User domain/database representation
type User struct {
	UUID      uuid.UUID      `gorm:"type:uuid;primaryKey" json:"uuid" binding:"-"`
	CreatedAt time.Time      `json:"createdAt" binding:"-"`
	UpdatedAt time.Time      `json:"updatedAt" binding:"-"`
	DeletedAt gorm.DeletedAt `json:"-" binding:"-"`
	FirstName string         `json:"firstName" binding:"required"`
	LastName  string         `json:"lastName"`
	Email     string         `json:"email" gorm:"UNIQUE" binding:"required"`
	Username  string         `json:"username" gorm:"UNIQUE" binding:"required"`
	Password  string         `json:"password"`
	Active    bool           `json:"active"  `
}

// TableName set user table name
func (User) TableName() string {
	return "public.users"
}

// BeforeCreate execute commands before creating User
func (u *User) BeforeCreate(scope *gorm.DB) (err error) {
	u.UUID = uuid.New()

	return
}
