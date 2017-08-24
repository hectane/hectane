package sql

import (
	"fmt"
	"strings"
)

const (
	OpEq = "="
)

// Clause represents selection criteria for the SQL SELECT statement. When
// generating the statement, the String() method is used to build the WHERE
// clause and the Values() method returns any parameters required.
type Clause interface {
	String() string
	Values() []interface{}
}

// ComparisonClause provides a comparison.
type ComparisonClause struct {
	Field    string
	Operator string
	Value    interface{}
}

func (c *ComparisonClause) String() string {
	return fmt.Sprintf("%s %s $", safeName(c.Field), c.Operator)
}

func (c *ComparisonClause) Values() []interface{} {
	return []interface{}{c.Value}
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
