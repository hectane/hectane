package util

import (
	"fmt"
	"reflect"
)

// DeleteItem composes an SQL DELETE statement for the specified value. v must
// be a pointer and the first field must be an ID.
func DeleteItem(t *Token, v interface{}) error {
	var (
		fnVal    = reflect.ValueOf(t.exec)
		fnParams = []reflect.Value{
			reflect.ValueOf(fmt.Sprintf(
				"DELETE FROM %s WHERE ID = $1",
				reflect.TypeOf(v).Elem().Name(),
			)),
			reflect.Indirect(reflect.ValueOf(v)).Field(0),
		}
	)
	if errVal := fnVal.Call(fnParams)[1]; !errVal.IsNil() {
		return errVal.Interface().(error)
	}
	return nil
}
