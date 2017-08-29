package resources

import (
	"fmt"
	"reflect"

	"github.com/hectane/hectane/db"
	"github.com/jinzhu/gorm"
	"github.com/manyminds/api2go"
)

// apply applies the query parameters to an SQL query. In order to ensure that
// the user cannot override fields, this method must be run before any other
// hooks so that they can override any fields.
func (r *Resource) apply(req api2go.Request) *gorm.DB {
	var (
		c        = db.C
		itemType = reflect.TypeOf(r.Type).Elem()
	)
	for i := 0; i < itemType.NumField(); i++ {
		n := gorm.ToDBName(itemType.Field(i).Name)
		v, ok := req.QueryParams[n]
		if ok {
			c = c.Where(fmt.Sprintf("%s = ?", n), v[0])
		}
	}
	return c
}
