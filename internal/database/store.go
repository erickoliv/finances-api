package database

import (
	"github.com/erickoliv/finances-api/accounts/accountsql"
	"github.com/erickoliv/finances-api/auth/sql"
	"github.com/erickoliv/finances-api/categories/categorysql"
	"github.com/erickoliv/finances-api/entries/entrysql"
	"github.com/erickoliv/finances-api/tags/tagsql"
	"gorm.io/gorm"
)

type SQLStore struct {
	Auth     *sql.AuthRepo
	Account  *accountsql.Repository
	Category *categorysql.CategoryRepo
	Tag      *tagsql.Repository
	Entry    *entrysql.Repository
}

func BuildSQLStore(db *gorm.DB) SQLStore {
	return SQLStore{
		Account:  accountsql.MakeAccounts(db),
		Auth:     sql.MakeAuthenticator(db),
		Category: categorysql.BuildRepository(db),
		Tag:      tagsql.BuildTagRepository(db),
		Entry:    entrysql.BuildRepository(db),
	}
}
