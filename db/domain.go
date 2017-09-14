package db

import (
	"strconv"
)

// Domain represents a FQDN used for routing incoming email and validating
// outgoing email.
type Domain struct {
	ID   int64  `json:"-"`
	Name string `json:"name" gorm:"not null;unique_index"`
}

func (d *Domain) GetName() string {
	return "domains"
}

func (d *Domain) GetID() string {
	return strconv.FormatInt(d.ID, 10)
}

func (d *Domain) SetID(id string) error {
	d.ID, _ = strconv.ParseInt(id, 10, 64)
	return nil
}
