package domain

const (
	// LoggedUser contains the UUID for the current user. It's set in AuthMidleware
	LoggedUser string = "current-logged-user" // maybe inject as parameter in the future ?
)
