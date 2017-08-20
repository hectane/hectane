package db

import (
	"github.com/hectane/hectane/db/util"
)

// Domain represents a FQDN used for routing incoming email and validating
// outgoing email.
type Domain struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func migrateDomainTable(t *util.Token) error {
	_, err := t.Exec(
		`
CREATE TABLE IF NOT EXISTS Domain (
	ID   SERIAL PRIMARY KEY,
	Name VARCHAR(80) NOT NULL UNIQUE
)
        `,
	)
	return err
}
