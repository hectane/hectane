package db

import (
	"reflect"
	"runtime"

	"github.com/hectane/hectane/db/util"
)

var migrations = []func(t *util.Token) error{
	migrateEntryTable,
	migrateUserTable,
	migrateDomainTable,
	migrateAccountTable,
	migrateFolderTable,
	migrateMessageTable,
}

// Migrate attempts to perform all pending database migrations. This function
// should be idempotent.
func Migrate() error {
	return Token.Transaction(func(t *util.Token) error {
		for _, m := range migrations {
			log.Debugf(
				"running \"%s\"...",
				runtime.FuncForPC(reflect.ValueOf(m).Pointer()).Name(),
			)
			if err := m(t); err != nil {
				return err
			}
		}
		return nil
	})
}
