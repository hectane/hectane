package receiver

import (
	"net/mail"
	"strings"

	"github.com/hectane/hectane/db"
	"github.com/jinzhu/gorm"
)

// lookup attempts to find an account that matches the specified email address.
func lookup(c *gorm.DB, address string) (*db.Account, error) {
	a, err := mail.ParseAddress(address)
	if err != nil {
		return nil, err
	}
	var (
		parts   = strings.Split(a.Address, "@")
		account = &db.Account{}
	)
	err = c.Table("accounts").
		Joins("LEFT JOIN domains ON accounts.domain_id = domains.id").
		Where("accounts.name = ?", parts[0]).
		Where("domains.name = ?", parts[1]).
		First(account).
		Error
	if err != nil {
		return nil, err
	}
	return account, nil
}
