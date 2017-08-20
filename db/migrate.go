package db

import (
	"github.com/hectane/hectane/db/util"
)

var migrations = []string{
	EntryModel,
	UserModel,
	DomainModel,
	AccountModel,
	FolderModel,
	MessageModel,
}

// Migrate attempts to perform all pending database migrations. This function
// should be idempotent.
func Migrate() error {
	return Token.Transaction(func(t *util.Token) error {
		for _, m := range migrations {
			log.Debugf(
				"migrating \"%s\" model...",
				m,
			)
			if err := modelRegistry[m].Migrate(t); err != nil {
				return err
			}
		}
		return nil
	})
}
