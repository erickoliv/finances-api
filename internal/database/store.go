package database

import (
	"github.com/erickoliv/finances-api/accounts/accountsql"
	"github.com/erickoliv/finances-api/entries/entrysql"
	"github.com/erickoliv/finances-api/repository/sql"
	"github.com/erickoliv/finances-api/tags/tagsql"
	"github.com/jinzhu/gorm"
)

type SQLStore struct {
	Auth     *sql.AuthRepo
	Account  *accountsql.Repository
	Category *sql.CategoryRepo // TODO: move this repository to categorysql
	Tag      *tagsql.Repository
	Entry    *entrysql.Repository
}

func BuildSQLStore(db *gorm.DB) SQLStore {
	return SQLStore{
		Account:  accountsql.MakeAccounts(db),
		Auth:     sql.MakeAuthenticator(db),
		Category: sql.BuildCategoryRepository(db),
		Tag:      tagsql.BuildTagRepository(db),
		Entry:    entrysql.BuildRepository(db),
	}
}
