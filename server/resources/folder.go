package resources

import (
	"github.com/hectane/hectane/db"
	"github.com/jinzhu/gorm"
	"github.com/manyminds/api2go"
)

const contextUser = "user"

var FolderResource = &Resource{
	Type: &db.Folder{},
	SetHook: func(obj interface{}, req api2go.Request) {
		u := req.PlainRequest.Context().Value(contextUser).(*db.User)
		obj.(*db.Folder).UserID = u.ID
	},
	GetHook: func(c *gorm.DB, req api2go.Request) *gorm.DB {
		u := req.PlainRequest.Context().Value(contextUser).(*db.User)
		return c.Where("user_id = ?", u.ID)
	},
}
