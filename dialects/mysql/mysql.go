// Package mysql implements reform.Dialect for MySQL.
package mysql // TODO add canonical import path via gopkg.in

import (
	"github.com/go-reform/reform"
)

type mysql struct{}

func (mysql) Placeholder(index int) string {
	return "?"
}

func (mysql) Placeholders(start, count int) []string {
	res := make([]string, count)
	for i := 0; i < count; i++ {
		res[i] = "?"
	}
	return res
}

func (mysql) QuoteIdentifier(identifier string) string {
	return "`" + identifier + "`"
}

func (mysql) LastInsertIdMethod() reform.LastInsertIdMethod {
	return reform.LastInsertId
}

// Dialect implements reform.Dialect for MySQL.
var Dialect mysql

// check interface
var _ reform.Dialect = Dialect
