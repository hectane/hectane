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
	).Error
	if err != nil {
		return err
	}
	// TODO: unique indexes
	return nil
}
