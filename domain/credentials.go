package domain

// Credentials is intended to be used to authentication operations
type Credentials struct {
	Username string `json:"username" binding:"required" `
	Password string `json:"password" binding:"required" `
}
