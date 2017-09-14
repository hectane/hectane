package resource

import (
	"net/http"

	"github.com/hectane/hectane/db"
	"github.com/jinzhu/gorm"
	"github.com/manyminds/api2go"
)

// Update attempts to update a model instance in the database.
func (r *Resource) Update(obj interface{}, req api2go.Request) (api2go.Responder, error) {
	if r.AllHook != nil {
		if err := r.AllHook(ActionUpdate, req); err != nil {
			return nil, err
		}
	}
	c := db.C
	if r.GetHook != nil {
		c = r.GetHook(c, req)
	}
	if r.SetHook != nil {
		r.SetHook(obj, req)
	}
	if err := db.Transaction(c, func(c *gorm.DB) error {
		if err := c.Model(r.Type).Updates(obj).Error; err != nil {
			return err
		}
		for _, p := range r.Preloads {
			c = c.Preload(p)
		}
		return c.First(obj).Error
	}); err != nil {
		return nil, err
	}
	return &api2go.Response{
		Res:  obj,
		Code: http.StatusOK,
	}, nil
}
