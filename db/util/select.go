package util

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

var (
	replacePlaceholders = regexp.MustCompile(`\$`)

	ErrRowCount = errors.New("invalid row count")
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
		fnVal      = reflect.ValueOf(t.Query)
		fnParams   []reflect.Value
	)
	for i := 0; i < itemType.NumField(); i++ {
		fieldNames = append(fieldNames, safeName(itemType.Field(i).Name))
	}
	query := fmt.Sprintf(
		"SELECT %s FROM %s",
		strings.Join(fieldNames, ", "),
		safeName(itemType.Name()),
	)
	if p.Where != nil {
		query += fmt.Sprintf(" WHERE %s", p.Where.String())
		for _, v := range p.Where.Values() {
			fnParams = append(fnParams, reflect.ValueOf(v))
		}
	}
	if p.OrderBy != "" {
		query += fmt.Sprintf(" ORDER BY %s", safeName(p.OrderBy))
		if p.OrderDesc {
			query += " DESC"
		}
	}
	i := 0
	query = replacePlaceholders.ReplaceAllStringFunc(query, func(string) string {
		i++
		return fmt.Sprintf("$%d", i)
	})
	fnParams = append([]reflect.Value{reflect.ValueOf(query)}, fnParams...)
	var (
		r      = fnVal.Call(fnParams)
		rVal   = r[0]
		errVal = r[1]
	)
	if !errVal.IsNil() {
		return nil, errVal.Interface().(error)
	}
	return rowsToItemSlice(rVal.Interface().(*sql.Rows), v)
}

// SelectItem is identical to SelectItems but ensures that the query returns
// only a single item. The first return value is the item, if successful.
func SelectItem(t *Token, v interface{}, p SelectParams) (interface{}, error) {
	i, err := SelectItems(t, v, p)
	if err != nil {
		return nil, err
	}
	itemVal := reflect.ValueOf(i)
	if itemVal.Len() != 1 {
		return nil, ErrRowCount
	}
	return itemVal.Index(0).Interface(), nil
}
