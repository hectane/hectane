package server

import (
	"errors"

	"github.com/hectane/hectane/db"
	"github.com/hectane/hectane/server/auth"
	"github.com/hectane/hectane/server/resource"
	"github.com/jinzhu/gorm"
	"github.com/manyminds/api2go"
)

var (
	ErrInsufficientPrivileges = errors.New("insufficient privileges")
	ErrInvalidAction          = errors.New("invalid action")
)

// requireAdmin ensures that only an administrator can perform the action.
func requireAdmin(_ int, req api2go.Request) error {
	u := req.PlainRequest.Context().Value(auth.User).(*db.User)
	if !u.IsAdmin {
		return ErrInsufficientPrivileges
	}
	return nil
}

// preventCreate prevents resources from being created.
func preventCreate(action int, req api2go.Request) error {
	if action == resource.ActionCreate {
		return ErrInvalidAction
	}
	return nil
}

var (
	domainResource = &resource.Resource{
		Type:    &db.Domain{},
		AllHook: requireAdmin,
	}
	folderResource = &resource.Resource{
		Type: &db.Folder{},
		SetHook: func(obj interface{}, req api2go.Request) {
			u := req.PlainRequest.Context().Value(auth.User).(*db.User)
			obj.(*db.Folder).UserID = u.ID
		},
		GetHook: func(c *gorm.DB, req api2go.Request) *gorm.DB {
			u := req.PlainRequest.Context().Value(auth.User).(*db.User)
			return c.Where("user_id = ?", u.ID).Order("name")
		},
	}
	messageResource = &resource.Resource{
		Type:    &db.Message{},
		Fields:  []string{"folder_id"},
		AllHook: preventCreate,
		SetHook: func(obj interface{}, req api2go.Request) {
			u := req.PlainRequest.Context().Value(auth.User).(*db.User)
			obj.(*db.Message).UserID = u.ID
		},
		GetHook: func(c *gorm.DB, req api2go.Request) *gorm.DB {
			u := req.PlainRequest.Context().Value(auth.User).(*db.User)
			return c.Where("user_id = ?", u.ID)
		},
	}
	userResource = &resource.Resource{
		Type:    &db.User{},
		AllHook: requireAdmin,
	}
)
