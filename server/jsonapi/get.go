package jsonapi

import (
	"net/url"
	"reflect"
	"strings"

	"github.com/hectane/hectane/db"
	"github.com/hectane/hectane/db/models"
	"github.com/hectane/hectane/db/sql"
)

func jsonToField(v interface{}, jsonName string) string {
	itemType := reflect.TypeOf(v)
	for i := 1; i < itemType.NumField(); i++ {
		var (
			field = itemType.Field(i)
			tag   = field.Tag.Get("json")
		)
		if idx := strings.Index(tag, ","); idx != -1 {
			tag = tag[:idx]
		}
		if tag == jsonName {
			return field.Name
		}
	}
	return ""
}

// Get retrieves items from the specified model, filtering the query using GET
// parameters.
func Get(model string, fields url.Values, p sql.SelectParams) (interface{}, error) {
	m, err := models.Get(model)
	if err != nil {
		return nil, err
	}
	clauses := sql.AndClause{}
	if p.Where != nil {
		clauses = append(clauses, p.Where)
	}
	for k, _ := range fields {
		if f := jsonToField(m.Instance, k); f != "" {
			clauses = append(clauses, &sql.ComparisonClause{
				Field:    f,
				Operator: sql.OpEq,
				Value:    fields.Get(k),
			})
		}
	}
	if len(clauses) != 0 {
		p.Where = &clauses
	}
	return sql.SelectItems(db.Token, m.Instance, p)
}
