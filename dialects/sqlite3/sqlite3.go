// Package sqlite3 implements reform.Dialect for SQLite3.
package sqlite3 // import "gopkg.in/reform.v1/dialects/sqlite3"

import (
	"gopkg.in/reform.v1"
)

type sqlite3 struct{}

func (sqlite3) String() string {
	return "sqlite3"
}

func (sqlite3) Placeholder(index int) string {
	return "?"
}

func (sqlite3) Placeholders(start, count int) []string {
	res := make([]string, count)
	for i := 0; i < count; i++ {
		res[i] = "?"
	}
	return res
}

func (sqlite3) QuoteIdentifier(identifier string) string {
	return `"` + identifier + `"`
}

func (sqlite3) LastInsertIdMethod() reform.LastInsertIdMethod {
	return reform.LastInsertId
}

func (sqlite3) SelectLimitMethod() reform.SelectLimitMethod {
	return reform.Limit
}

func (sqlite3) DefaultValuesMethod() reform.DefaultValuesMethod {
	return reform.DefaultValues
}

// Dialect implements reform.Dialect for SQLite3.
var Dialect sqlite3

// check interface
var _ reform.Dialect = Dialect
