package resource

import (
	"net/http"
	"reflect"

	"github.com/manyminds/api2go"
)

// FindAll attempts to retrieve all instances of a model from the database.
func (r *Resource) FindAll(req api2go.Request) (api2go.Responder, error) {
	if r.AllHook != nil {
		if err := r.AllHook(ActionFindAll, req); err != nil {
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
		itemType  = reflect.TypeOf(r.Type)
		sliceType = reflect.SliceOf(itemType)
		sliceVal  = reflect.New(sliceType)
	)
	if err := c.Find(sliceVal.Interface()).Error; err != nil {
		return nil, err
	}
	return &api2go.Response{
		Res:  reflect.Indirect(sliceVal).Interface(),
		Code: http.StatusOK,
	}, nil
}
