package db

import (
	"database/sql"
	"fmt"

	dbsql "github.com/hectane/hectane/db/sql"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

var (
	Token *dbsql.Token
	log   = logrus.WithField("context", "db")
)

// Connect establishes a connection to the PostgreSQL database used for all SQL
// queries. This function should be called before using any other types or
// functions in the package.
func Connect(name, user, password, host string, port int) error {
	c, err := sql.Open(
		"postgres",
		fmt.Sprintf(
			"dbname=%s user=%s password=%s host=%s port=%d",
			name,
			user,
			password,
			host,
			port,
		),
	)
	if err != nil {
		return err
	}
	Token = dbsql.NewToken(c)
	return nil
}
