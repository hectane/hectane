package resource

import (
	"net/http"
	"reflect"

	"github.com/manyminds/api2go"
)

// FindOne attempts to retrieve a single model instance from the database.
func (r *Resource) FindOne(ID string, req api2go.Request) (api2go.Responder, error) {
	if r.AllHook != nil {
		if err := r.AllHook(ActionFindOne, req); err != nil {
			return nil, err
		}
	}
	c, err := r.apply(req)
	if err != nil {
		return nil, err
	}
	if r.GetHook != nil {
		c = r.GetHook(c, req)
	}
	var (
		itemType = reflect.TypeOf(r.Type).Elem()
		itemVal  = reflect.New(itemType)
	)
	if err := c.First(itemVal.Interface(), ID).Error; err != nil {
		return nil, err
	}
	return &api2go.Response{
		Res:  itemVal.Interface(),
		Code: http.StatusOK,
	}, nil
}
