package resources

import (
	"errors"
	"fmt"

	"github.com/hectane/hectane/db"
	"github.com/jinzhu/gorm"
	"github.com/manyminds/api2go"
)

var ErrInvalidParameter = errors.New("invalid parameter")

// apply applies the query parameters to an SQL query. In order to ensure that
// the user cannot override fields, this method must be run before any other
// hooks so that they can override any fields.
func (r *Resource) apply(req api2go.Request) (*gorm.DB, error) {
	c := db.C
loop:
	for k, v := range req.QueryParams {
		for _, f := range r.Fields {
			if f == k {
				c = c.Where(fmt.Sprintf("%s = ?", k), v[0])
				continue loop
			}
		}
		return nil, ErrInvalidParameter
	}
	return c, nil
}
