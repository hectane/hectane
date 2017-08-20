package db

import (
	"github.com/hectane/hectane/db/util"
)

const (
	EntryModel   = "entry"
	UserModel    = "user"
	DomainModel  = "domain"
	AccountModel = "account"
	FolderModel  = "folder"
	MessageModel = "message"
)

// Model contains metadata necessary for interacting with data in a specific
// database table.
type Model struct {
	Instance interface{}
	Migrate  func(t *util.Token) error
}

var modelRegistry = map[string]Model{
	EntryModel: Model{
		Instance: Entry{},
		Migrate:  migrateEntryTable,
	},
	UserModel: Model{
		Instance: User{},
		Migrate:  migrateUserTable,
	},
	DomainModel: Model{
		Instance: Domain{},
		Migrate:  migrateDomainTable,
	},
	AccountModel: Model{
		Instance: Account{},
		Migrate:  migrateAccountTable,
	},
	FolderModel: Model{
		Instance: Folder{},
		Migrate:  migrateFolderTable,
	},
	MessageModel: Model{
		Instance: Message{},
		Migrate:  migrateMessageTable,
	},
}
