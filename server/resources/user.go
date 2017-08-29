package resources

import (
	"github.com/hectane/hectane/db"
)

var UserResource = &Resource{
	Type:    &db.User{},
	AllHook: requireAdmin,
}
