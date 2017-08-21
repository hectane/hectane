package models

import (
	"errors"

	"github.com/hectane/hectane/db/sql"
)

var (
	ErrInvalidModel = errors.New("invalid model specified")
)

// Model contains metadata necessary for interacting with data in a specific
// database table.
type Model struct {
	Instance interface{}
	Migrate  func(t *sql.Token) error
}

var registry = map[string]*Model{
	Entries: &Model{
		Instance: Entry{},
		Migrate:  migrateEntries,
	},
	Users: &Model{
		Instance: User{},
		Migrate:  migrateUsers,
	},
	Domains: &Model{
		Instance: Domain{},
		Migrate:  migrateDomains,
	},
	Accounts: &Model{
		Instance: Account{},
		Migrate:  migrateAccounts,
	},
	Folders: &Model{
		Instance: Folder{},
		Migrate:  migrateFolders,
	},
	Messages: &Model{
		Instance: Message{},
		Migrate:  migrateMessages,
	},
}

// Get retrieves the metadata for the specified model if available.
func Get(name string) (*Model, error) {
	m, ok := registry[name]
	if !ok {
		return nil, ErrInvalidModel
	}
	return m, nil
}
