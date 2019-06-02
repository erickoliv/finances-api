package model

// User model/database representation
type User struct {
	BaseModel
	FirstName string `json:"firstName" binding:"required"`
	LastName  string `json:"lastName"`
	Email     string `json:"email" gorm:"UNIQUE" `
	Username  string `json:"username" gorm:"UNIQUE" `
	Password  string `json:"password"  `
	Active    bool   `json:"active"  `
}

// TableName set user table name
func (User) TableName() string {
	return "public.users"
}
