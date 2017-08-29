package resources

import (
	"net/http"
	"reflect"

	"github.com/hectane/hectane/db"
	"github.com/jinzhu/gorm"
	"github.com/manyminds/api2go"
)

// Resource implements the interfaces necessary to use a database model with
// the api2go package. Hooks can be used to apply filtering to the methods.
type Resource struct {
	Type    interface{}
	SetHook func(interface{}, api2go.Request)
	GetHook func(*gorm.DB, api2go.Request) *gorm.DB
}

// Create attempts to save a new model instance to the database.
func (r *Resource) Create(obj interface{}, req api2go.Request) (api2go.Responder, error) {
	if r.SetHook != nil {
		r.SetHook(obj, req)
	}
	if err := db.C.Create(obj).Error; err != nil {
		return nil, err
	}
	return &api2go.Response{
		Res:  obj,
		Code: http.StatusCreated,
	}, nil
}

// Delete attempts to delete the specified model instance from the database.
func (r *Resource) Delete(id string, req api2go.Request) (api2go.Responder, error) {
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

// FindAll attempts to retrieve all instances of a model from the database.
func (r *Resource) FindAll(req api2go.Request) (api2go.Responder, error) {
	c := db.C
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
		Code: 200,
	}, nil
}

// FindOne attempts to retrieve a single model instance from the database.
func (r *Resource) FindOne(ID string, req api2go.Request) (api2go.Responder, error) {
	c := db.C
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
		Code: 200,
	}, nil
}

// Update attempts to update a model instance in the database.
func (r *Resource) Update(obj interface{}, req api2go.Request) (api2go.Responder, error) {
	c := db.C
	if r.GetHook != nil {
		c = r.GetHook(c, req)
	}
	if r.SetHook != nil {
		r.SetHook(obj, req)
	}
	if err := c.Model(r.Type).Updates(obj).Error; err != nil {
		return nil, err
	}
	return &api2go.Response{
		Res:  obj,
		Code: http.StatusOK,
	}, nil
}
