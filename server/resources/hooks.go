package resources

import (
	"errors"

	"github.com/hectane/hectane/db"
	"github.com/manyminds/api2go"
)

var (
	ErrInsufficientPrivileges = errors.New("insufficient privileges")
	ErrInvalidAction          = errors.New("invalid action")
)

// requireAdmin ensures that only an administrator can perform the action.
func requireAdmin(_ int, req api2go.Request) error {
	u := req.PlainRequest.Context().Value(contextUser).(*db.User)
	if !u.IsAdmin {
		return ErrInsufficientPrivileges
	}
	return nil
}

// preventCreate prevents resources from being created.
func preventCreate(action int, req api2go.Request) error {
	if action == actionCreate {
		return ErrInvalidAction
	}
	return nil
}
