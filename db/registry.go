package db

import (
	"github.com/hectane/hectane/db/models"
	"github.com/hectane/hectane/db/util"
)

// Model contains metadata necessary for interacting with data in a specific
// database table.
type Model struct {
	Instance interface{}
	Migrate  func(t *util.Token) error
}

var modelRegistry = map[string]Model{
	models.Entries: Model{
		Instance: models.Entry{},
		Migrate:  models.MigrateEntries,
	},
	models.Users: Model{
		Instance: models.User{},
		Migrate:  models.MigrateUsers,
	},
	models.Domains: Model{
		Instance: models.Domain{},
		Migrate:  models.MigrateDomains,
	},
	models.Accounts: Model{
		Instance: models.Account{},
		Migrate:  models.MigrateAccounts,
	},
	models.Folders: Model{
		Instance: models.Folder{},
		Migrate:  models.MigrateFolders,
	},
	models.Messages: Model{
		Instance: models.Message{},
		Migrate:  models.MigrateMessages,
	},
}
