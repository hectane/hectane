package db

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/sirupsen/logrus"
)

var (
	C   *gorm.DB
	log = logrus.WithField("context", "db")
)

// Connect establishes a connection to the database.
func Connect(driver, name, user, password, host string, port int) error {
	d, err := gorm.Open(
		driver,
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
	C = d
	return nil
}

// Migrate performs all database migrations.
func Migrate() error {
	log.Info("performing migrations...")
	err := C.AutoMigrate(
		&User{},
		&Domain{},
		&Account{},
		&Folder{},
		&Message{},
		&Contact{},
		&QueueItem{},
	).Error
	if err != nil {
		return err
	}
	// TODO: unique indexes
	return nil
}

// Transaction begins a new transaction which will either rollback or commit
// based on the return value of the callback.
func Transaction(c *gorm.DB, fn func(*gorm.DB) error) error {
	c = c.Begin()
	if err := fn(c); err != nil {
		c.Rollback()
		return err
	}
	c.Commit()
	return nil
}
