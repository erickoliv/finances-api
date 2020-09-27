package database

import (
	"strings"

	"github.com/erickoliv/finances-api/accounts"
	"github.com/erickoliv/finances-api/auth"
	"github.com/erickoliv/finances-api/categories"
	"github.com/erickoliv/finances-api/entries"
	"github.com/erickoliv/finances-api/tags"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	migrations := []interface{}{
		tags.Tag{},
		auth.User{},
		categories.Category{},
		accounts.Account{},
		entries.Entry{},
		entries.EntryTag{},
	}

	for _, entity := range migrations {
		err := db.AutoMigrate(entity)
		if err == nil {
			continue
		}

		// ignoring these errors for now
		if strings.Contains(err.Error(), "already exists") {
			continue
		}

		return err
	}

	return db.Error
}
