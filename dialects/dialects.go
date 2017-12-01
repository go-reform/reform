// Package dialects implements reform.Dialect selector.
package dialects

import (
	"gopkg.in/reform.v1"
	"gopkg.in/reform.v1/dialects/mssql"
	"gopkg.in/reform.v1/dialects/mysql"
	"gopkg.in/reform.v1/dialects/postgresql"
	"gopkg.in/reform.v1/dialects/sqlite3"
	"gopkg.in/reform.v1/dialects/sqlserver"
)

// ForDriver returns reform Dialect for given driver string, or nil.
func ForDriver(driver string) reform.Dialect {
	switch driver {
	case "postgres":
		return postgresql.Dialect
	case "mysql":
		return mysql.Dialect
	case "sqlite3":
		return sqlite3.Dialect
	case "mssql":
		return mssql.Dialect
	case "sqlserver":
		return sqlserver.Dialect
	default:
		return nil
	}
}
