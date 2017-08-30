package resource

import (
	"net/http"

	"github.com/hectane/hectane/db"
	"github.com/manyminds/api2go"
)

// Delete attempts to delete the specified model instance from the database.
func (r *Resource) Delete(id string, req api2go.Request) (api2go.Responder, error) {
	if r.AllHook != nil {
		if err := r.AllHook(ActionDelete, req); err != nil {
			return nil, err
		}
	}
	c := db.C
	if r.GetHook != nil {
		c = r.GetHook(c, req)
	}
	if err := c.Where("ID = ?", id).Delete(r.Type).Error; err != nil {
		return nil, err
	}
	return &api2go.Response{
		Code: http.StatusOK,
	}, nil
}
