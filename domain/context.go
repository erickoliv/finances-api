package domain

const (
	// DB is used to pass database connection reference inside gin.Context
	DB string = "DB_POOL"
	// AppToken is a string used to hash passwords, extracted from environment variable
	AppToken string = "APP_TOKEN"
	// LoggedUser contains the UUID for the current user. It's set in AuthMidleware
	LoggedUser string = "current-logged-user" // maybe inject as parameter in the future ?
)
