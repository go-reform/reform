// Package sqlite3 implements reform.Dialect for SQLite3.
package sqlite3 // TODO add canonical import path via gopkg.in

import (
	"github.com/AlekSi/reform"
)

type sqlite3 struct{}

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

// Dialect implements reform.Dialect for SQLite3.
var Dialect sqlite3

// check interface
var _ reform.Dialect = Dialect
