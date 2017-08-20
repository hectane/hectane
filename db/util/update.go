package util

import (
	"fmt"
	"reflect"
	"strings"
)

// UpdateItem composes an SQL UPDATE statement for the specified value, setting
// all of the fields to their new value from the struct. v must be a pointer
// and the first field must be an ID.
func UpdateItem(t *Token, v interface{}) error {
	var (
		itemType         = reflect.TypeOf(v).Elem()
		itemVal          = reflect.Indirect(reflect.ValueOf(v))
		fieldAssignments []string
		fnVal            = reflect.ValueOf(t.Exec)
		fnParams         []reflect.Value
	)
	for i := 1; i < itemType.NumField(); i++ {
		fieldAssignments = append(
			fieldAssignments,
			fmt.Sprintf("%s=$%d", safeName(itemType.Field(i).Name), i),
		)
		fnParams = append(fnParams, itemVal.Field(i))
	}
	fnParams = append(fnParams, itemVal.Field(0))
	fnParams = append([]reflect.Value{reflect.ValueOf(fmt.Sprintf(
		"UPDATE %s SET %s WHERE ID = $%d",
		safeName(itemType.Name()),
		strings.Join(fieldAssignments, ", "),
		itemType.NumField(),
	))}, fnParams...)
	if errVal := fnVal.Call(fnParams)[1]; !errVal.IsNil() {
		return errVal.Interface().(error)
	}
	return nil
}
