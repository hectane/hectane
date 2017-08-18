package db

import (
	"database/sql"
	"fmt"
)

// Account represents an individual email account owned by a user.
type Account struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	UserID   int    `json:"user_id"`
	DomainID int    `json:"domain_id"`
}

func migrateAccountsTable(t *Token) error {
	_, err := t.exec(
		`
CREATE TABLE IF NOT EXISTS Accounts (
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

func rowsToAccounts(r *sql.Rows) ([]*Account, error) {
	accounts := make([]*Account, 0, 1)
	for r.Next() {
		a := &Account{}
		if err := r.Scan(&a.ID, &a.Name, &a.UserID, &a.DomainID); err != nil {
			return nil, err
		}
		accounts = append(accounts, a)
	}
	return accounts, nil
}

// Accounts retrieves a list of all accounts in the database.
func Accounts(t *Token) ([]*Account, error) {
	r, err := t.query(
		`
SELECT ID, Name, UserID, DomainID
FROM Users ORDER BY Name
        `,
	)
	if err != nil {
		return nil, err
	}
	return rowsToAccounts(r)
}

// FindAccount attempts to find an account where the specified field matches
// the specified value. Exactly one row must be returned.
func FindAccount(t *Token, field string, value interface{}) (*Account, error) {
	a, err := FindAccounts(t, field, value)
	if err != nil {
		return nil, err
	}
	if len(a) != 1 {
		return nil, ErrRowCount
	}
	return a[0], nil
}

// FindAccounts attempts to retrieve all accounts where the specified field
// matches the specified value.
func FindAccounts(t *Token, field string, value interface{}) ([]*Account, error) {
	r, err := t.query(
		fmt.Sprintf(
			`
SELECT ID, Name, UserID, DomainID
FROM Accounts WHERE %s = $1 ORDER BY Name
            `,
			field,
		),
		value,
	)
	if err != nil {
		return nil, err
	}
	return rowsToAccounts(r)
}

// Save persists changes to the account. If ID is set to zero, a new account is
// created and its ID updated.
func (a *Account) Save(t *Token) error {
	if a.ID == 0 {
		err := t.queryRow(
			`
INSERT INTO Accounts (Name, UserID, DomainID)
VALUES ($1, $2, $3) RETURNING ID
            `,
			a.Name,
			a.UserID,
			a.DomainID,
		).Scan(&a.ID)
		if err != nil {
			return err
		}
		return nil
	} else {
		_, err := t.exec(
			`
UPDATE Accouts SET Name=$1, UserID=$2, DomainID=$3
WHERE ID = $4
            `,
			a.Name,
			a.UserID,
			a.DomainID,
			a.ID,
		)
		return err
	}
}

// Delete the account from the database.
func (a *Account) Delete(t *Token) error {
	_, err := t.exec(
		`
DELETE FROM Accounts WHERE ID = $1
        `,
		a.ID,
	)
	return err
}
