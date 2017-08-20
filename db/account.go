package db

import (
	"github.com/hectane/hectane/db/util"
)

// Account represents an individual email account owned by a user.
type Account struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	UserID   int    `json:"user_id"`
	DomainID int    `json:"domain_id"`
}

func migrateAccountTable(t *util.Token) error {
	_, err := t.Exec(
		`
CREATE TABLE IF NOT EXISTS Account (
    ID       SERIAL PRIMARY KEY,
    Name     VARCHAR(40) NOT NULL,
    UserID   INTEGER REFERENCES Users (ID) ON DELETE CASCADE,
    DomainID INTEGER REFERENCES Domains (ID) ON DELETE CASCADE,
    UNIQUE (Name, DomainID)
)
        `,
	)
	return err
}
