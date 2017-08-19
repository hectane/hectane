package db

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
)

// rowsToItemSlice converts a row of items to a slice of the same type as v. It
// is assumed that each row contains all of the fields in the order they
// appear.
func rowsToItemSlice(v interface{}, r *sql.Rows) (interface{}, error) {
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

// insertItem composes an SQL INSERT statement for the specified value, setting
// all of the fields to their equivalent value from the struct. v must be a
// pointer and the first field must be ID.
func insertItem(t *Token, v interface{}) error {
	var (
		itemType      = reflect.TypeOf(v).Elem()
		itemVal       = reflect.Indirect(reflect.ValueOf(v))
		fieldNames    []string
		fieldValues   []string
		fnQueryVal    = reflect.ValueOf(t.queryRow)
		fnQueryParams []reflect.Value
	)
	for i := 1; i < itemType.NumField(); i++ {
		fieldNames = append(fieldNames, itemType.Field(i).Name)
		fieldValues = append(fieldValues, fmt.Sprintf("$%d", i))
		fnQueryParams = append(fnQueryParams, itemVal.Field(i))
	}
	fnQueryParams = append([]reflect.Value{reflect.ValueOf(fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (%s) RETURNING ID",
		itemType.Name(),
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

// updateItem composes an SQL UPDATE statement for the specified value, setting
// all of the fields to their new value from the struct. v must be a pointer
// and the first field must be an ID.
func updateItem(t *Token, v interface{}) error {
	var (
		itemType         = reflect.TypeOf(v).Elem()
		itemVal          = reflect.Indirect(reflect.ValueOf(v))
		fieldAssignments []string
		fnVal            = reflect.ValueOf(t.exec)
		fnParams         []reflect.Value
	)
	for i := 1; i < itemType.NumField(); i++ {
		fieldAssignments = append(
			fieldAssignments,
			fmt.Sprintf("%s=$%d", itemType.Field(i).Name, i),
		)
		fnParams = append(fnParams, itemVal.Field(i))
	}
	fnParams = append(fnParams, itemVal.Field(0))
	fnParams = append([]reflect.Value{reflect.ValueOf(fmt.Sprintf(
		"UPDATE %s SET %s WHERE ID = $%d",
		itemType.Name(),
		strings.Join(fieldAssignments, ", "),
		itemType.NumField(),
	))}, fnParams...)
	if errVal := fnVal.Call(fnParams)[1]; !errVal.IsNil() {
		return errVal.Interface().(error)
	}
	return nil
}

// deleteItem composes an SQL DELETE statement for the specified value. v must
// be a pointer and the first field must be an ID.
func deleteItem(t *Token, v interface{}) error {
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
