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
		for _, n := range migrations {
			m, err := models.Get(n)
			if err != nil {
				return err
			}
			log.Debugf("migrating %s...", n)
			if err := m.Migrate(t); err != nil {
				return err
			}
		}
		return nil
	})
}
