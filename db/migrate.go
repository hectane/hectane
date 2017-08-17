package db

var migrations = []func(t *Token) error{
	migrateUsersTable,
	migrateDomainsTable,
	migrateAccountsTable,
	migrateFoldersTable,
}

// Migrate attempts to perform all pending database migrations. This function
// should be idempotent.
func Migrate() error {
	return Transaction(func(t *Token) error {
		for _, m := range migrations {
			if err := m(t); err != nil {
				return err
			}
		}
		return nil
	})
}
