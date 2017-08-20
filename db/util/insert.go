package util

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
)

// InsertItem composes an SQL INSERT statement for the specified value, setting
// all of the fields to their equivalent value from the struct. v must be a
// pointer and the first field must be ID.
func InsertItem(t *Token, v interface{}) error {
	var (
		itemType      = reflect.TypeOf(v).Elem()
		itemVal       = reflect.Indirect(reflect.ValueOf(v))
		fieldNames    []string
		fieldValues   []string
		fnQueryVal    = reflect.ValueOf(t.QueryRow)
		fnQueryParams []reflect.Value
	)
	for i := 1; i < itemType.NumField(); i++ {
		fieldNames = append(fieldNames, safeName(itemType.Field(i).Name))
		fieldValues = append(fieldValues, fmt.Sprintf("$%d", i))
		fnQueryParams = append(fnQueryParams, itemVal.Field(i))
	}
	fnQueryParams = append([]reflect.Value{reflect.ValueOf(fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (%s) RETURNING ID",
		safeName(itemType.Name()),
		strings.Join(fieldNames, ", "),
		strings.Join(fieldValues, ", "),
	))}, fnQueryParams...)
	var (
		r            = fnQueryVal.Call(fnQueryParams)[0].Interface().(*sql.Row)
		fnScanVal    = reflect.ValueOf(r.Scan)
		fnScanParams = []reflect.Value{itemVal.Field(0).Addr()}
	)
	if errVal := fnScanVal.Call(fnScanParams)[0]; !errVal.IsNil() {
		return errVal.Interface().(error)
	}
	return nil
}
