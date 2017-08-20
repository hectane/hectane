package util

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
)

// SelectParams provide parameters to the SQL SELECT statement.
type SelectParams struct {
	Where     Clause
	OrderBy   string
	OrderDesc bool
}

// SelectItems composes an SQL SELECT statement. All fields from v are
// retrieved and the parameters are used to filter the results.
func SelectItems(t *Token, v interface{}, p SelectParams) (interface{}, error) {
	var (
		itemType   = reflect.TypeOf(v)
		fieldNames []string
		fnVal      = reflect.ValueOf(t.query)
		fnParams   []reflect.Value
	)
	for i := 0; i < itemType.NumField(); i++ {
		fieldNames = append(fieldNames, itemType.Field(i).Name)
	}
	query := fmt.Sprintf(
		"SELECT %s FROM %s",
		strings.Join(fieldNames, ", "),
		itemType.Name(),
	)
	if p.Where != nil {
		query += fmt.Sprintf(" WHERE %s", p.Where.String())
		for _, v := range p.Where.Values() {
			fnParams = append(fnParams, reflect.ValueOf(v))
		}
	}
	if p.OrderBy != "" {
		query += fmt.Sprintf(" ORDER BY %s", p.OrderBy)
		if p.OrderDesc {
			query += " DESC"
		}
	}
	fnParams = append([]reflect.Value{reflect.ValueOf(query)}, fnParams...)
	r := fnVal.Call(fnParams)[0].Interface().(*sql.Rows)
	return rowsToItemSlice(r, v)
}
