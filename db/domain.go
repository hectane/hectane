package db

import (
	"github.com/jinzhu/gorm"
)

// Domain represents a FQDN used for routing incoming email and validating
// outgoing email.
type Domain struct {
	gorm.Model
	Name string `gorm:"not null;unique_index"`
}
