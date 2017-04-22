// Package sqlserver implements reform.Dialect for Microsoft SQL Server (sqlserver driver).
package sqlserver // import "gopkg.in/reform.v1/dialects/sqlserver"

import (
	"strconv"

	"gopkg.in/reform.v1"
)

type sqlserver struct{}

func (sqlserver) String() string {
	return "sqlserver"
}

func (sqlserver) Placeholder(index int) string {
	return "@P" + strconv.Itoa(index)
}

func (sqlserver) Placeholders(start, count int) []string {
	res := make([]string, count)
	for i := 0; i < count; i++ {
		res[i] = "@P" + strconv.Itoa(1+i)
	}
	return res
}

func (sqlserver) QuoteIdentifier(identifier string) string {
	return "[" + identifier + "]"
}

func (sqlserver) LastInsertIdMethod() reform.LastInsertIdMethod {
	return reform.OutputInserted
}

func (sqlserver) SelectLimitMethod() reform.SelectLimitMethod {
	return reform.SelectTop
}

func (sqlserver) DefaultValuesMethod() reform.DefaultValuesMethod {
	return reform.DefaultValues
}

// Dialect implements reform.Dialect for Microsoft SQL Server.
var Dialect sqlserver

// check interface
var _ reform.Dialect = Dialect
