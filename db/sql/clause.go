package sql

import (
	"fmt"
	"strings"
)

// Clause represents selection criteria for the SQL SELECT statement. When
// generating the statement, the String() method is used to build the WHERE
// clause and the Values() method returns any parameters required.
type Clause interface {
	String() string
	Values() []interface{}
}

// EqClause provides an extremely simple equality comparison.
type EqClause struct {
	Field string
	Value interface{}
}

func (e *EqClause) String() string {
	return fmt.Sprintf("%s = $", safeName(e.Field))
}

func (e *EqClause) Values() []interface{} {
	return []interface{}{e.Value}
}

// AndClause combines two or more clauses with the AND operator.
type AndClause []Clause

func (a *AndClause) String() string {
	var clauses []string
	for _, c := range *a {
		clauses = append(clauses, c.String())
	}
	return fmt.Sprintf("(%s)", strings.Join(clauses, " AND "))
}

func (a *AndClause) Values() []interface{} {
	var clauses []interface{}
	for _, c := range *a {
		clauses = append(clauses, c.Values()...)
	}
	return clauses
}
