package internal

import (
	"gopkg.in/reform.v1"
	"gopkg.in/reform.v1/dialects/mssql"
	"gopkg.in/reform.v1/dialects/mysql"
	"gopkg.in/reform.v1/dialects/postgresql"
	"gopkg.in/reform.v1/dialects/sqlite3"
)

// FIXME decide if we really want to introduce a new package.
// It's internal, but should be vendored, etc.

// DialectForDriver returns reform Dialect for given driver string.
func DialectForDriver(driver string) reform.Dialect {
	switch driver {
	case "postgres":
		return postgresql.Dialect
	case "mysql":
		return mysql.Dialect
	case "sqlite3":
		return sqlite3.Dialect
	case "mssql":
		return mssql.Dialect
	default:
		return nil
	}
}
