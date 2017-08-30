package resources

import (
	"github.com/hectane/hectane/db"
)

var DomainResource = &Resource{
	Type:    &db.Domain{},
	AllHook: requireAdmin,
}
