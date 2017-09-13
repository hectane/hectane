package resource

import (
	"errors"
	"fmt"

	"github.com/hectane/hectane/db"
	"github.com/jinzhu/gorm"
	"github.com/manyminds/api2go"
)

var ErrInvalidParameter = errors.New("invalid parameter")

// apply applies the query parameters to an SQL query.
func (r *Resource) apply(req api2go.Request) (*gorm.DB, error) {
	c := db.C
	for _, p := range r.Preloads {
		c = c.Preload(p)
	}
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
