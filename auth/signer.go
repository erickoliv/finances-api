package auth

type SessionSigner interface {
	SignUser(identifier string) (string, error)
	Validate(token string) (string, error)
}
