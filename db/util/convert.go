package util

import (
	"database/sql"
	"reflect"
)

// rowsToItemSlice converts a row of items to a slice of the same type as v. It
// is assumed that each row contains all of the fields in the order they
// appear.
func rowsToItemSlice(r *sql.Rows, v interface{}) (interface{}, error) {
	var (
		itemType  = reflect.TypeOf(v)
		sliceType = reflect.SliceOf(reflect.PtrTo(itemType))
		sliceVal  = reflect.MakeSlice(sliceType, 0, 1)
		fnVal     = reflect.ValueOf(r.Scan)
	)
	for r.Next() {
		var (
			itemVal  = reflect.New(itemType)
			fnParams = []reflect.Value{}
		)
		for i := 0; i < itemType.NumField(); i++ {
			fnParams = append(fnParams, reflect.Indirect(itemVal).Field(i).Addr())
		}
		if errVal := fnVal.Call(fnParams)[0]; !errVal.IsNil() {
			return nil, errVal.Interface().(error)
		}
		sliceVal = reflect.Append(sliceVal, itemVal)
	}
	return sliceVal.Interface(), nil
}
