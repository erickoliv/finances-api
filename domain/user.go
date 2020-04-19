package domain

// User domain/database representation
type User struct {
	BaseModel
	FirstName string `json:"firstName" binding:"required"`
	LastName  string `json:"lastName"`
	Email     string `json:"email" gorm:"UNIQUE" binding:"required"`
	Username  string `json:"username" gorm:"UNIQUE" binding:"required"`
	Password  string `json:"password"`
	Active    bool   `json:"active"  `
}

// TableName set user table name
func (User) TableName() string {
	return "public.users"
}
