package db

import (
	"github.com/hectane/hectane/db/models"
	"github.com/hectane/hectane/db/sql"
)

var migrations = []string{
	models.Entries,
	models.Users,
	models.Domains,
	models.Accounts,
	models.Folders,
	models.Messages,
}

// Migrate attempts to perform all pending database migrations. This function
// should be idempotent.
func Migrate() error {
	return Token.Transaction(func(t *sql.Token) error {
		for _, m := range migrations {
			log.Debugf("migrating %s...", m)
			if err := modelRegistry[m].Migrate(t); err != nil {
				return err
			}
		}
		return nil
	})
}
