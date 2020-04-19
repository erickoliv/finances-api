package sql

import (
	"context"
	"errors"

	"github.com/erickoliv/finances-api/domain"
	"github.com/erickoliv/finances-api/service"
	"github.com/jinzhu/gorm"
)

type authRepo struct {
	db *gorm.DB
}

var (
	errEmptyUsername = errors.New("username cannot be empty")
	errEmptyPassword = errors.New("password cannot be empty")
	errInvalidUser   = errors.New("user or password invalid")
)

// MakeAuthenticator returns a service.Authenticator sql implementation using gorm
func MakeAuthenticator(db *gorm.DB) service.Authenticator {
	return &authRepo{
		db,
	}
}

func (repo *authRepo) Login(ctx context.Context, username string, password string) (*domain.User, error) {
	if username == "" {
		return nil, errEmptyUsername
	}
	if password == "" {
		return nil, errEmptyPassword
	}

	user := &domain.User{}
	result := repo.db.First(user, "username = ? AND password = crypt(?, password)", username, password)
	if result.RecordNotFound() {
		return nil, errInvalidUser
	}

	return user, result.Error
}

func (repo *authRepo) Register(ctx context.Context, user *domain.User) error {
	if user == nil {
		return errInvalidUser
	}
	result := repo.db.Create(user)

	// TODO: move to a specific method
	repo.db.Model(user).Update("password", gorm.Expr("crypt(?, gen_salt('bf', 8))", user.Password))

	return result.Error
}
