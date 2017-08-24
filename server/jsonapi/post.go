package jsonapi

import (
	"encoding/json"
	"io"
	"reflect"

	"github.com/hectane/hectane/db/models"
)

// Post attempts to populate a new model instance for the specified model. The
// instance is not saved to the database.
func Post(model string, r io.ReadCloser) (interface{}, error) {
	m, err := models.Get(model)
	if err != nil {
		return nil, err
	}
	i := reflect.New(reflect.TypeOf(m.Instance)).Interface()
	if err := json.NewDecoder(r).Decode(i); err != nil {
		return nil, err
	}
	return i, err
}
