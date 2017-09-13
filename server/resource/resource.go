package resource

import (
	"github.com/jinzhu/gorm"
	"github.com/manyminds/api2go"
)

const (
	ActionCreate = iota
	ActionDelete
	ActionFindAll
	ActionFindOne
	ActionUpdate
)

// Resource implements the interfaces necessary to use a database model with
// the api2go package. Preloads determines which fields should be preloaded.
// Fields determines which fields can be used for filtering. Hooks can be used
// to apply filtering to the methods.
type Resource struct {
	Type     interface{}
	Preloads []string
	Fields   []string
	AllHook  func(int, api2go.Request) error
	SetHook  func(interface{}, api2go.Request)
	GetHook  func(*gorm.DB, api2go.Request) *gorm.DB
}
