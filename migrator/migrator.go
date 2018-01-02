package migrator

import (
	"fmt"
	"strings"

	reform "gopkg.in/reform.v1"
	"gopkg.in/reform.v1/dialects/mssql"
	"gopkg.in/reform.v1/dialects/mysql"
	"gopkg.in/reform.v1/dialects/postgresql"
	"gopkg.in/reform.v1/dialects/sqlite3"
	"gopkg.in/reform.v1/dialects/sqlserver"
)

const migrationTable = "schema_migrations"

// Manager is a migrator which stores information
// about upgraded/downgraded migrations in the database.
type Manager struct {
	db *reform.DB
}

func (m *Manager) createMigrationTable() error {
	type column struct {
		column     string
		datatype   string
		attributes string
	}

	ver := column{"version", "VARCHAR(50)", "PRIMARY KEY"}
	state := column{"state", "VARCHAR(10)", "NOT NULL"}
	runat := column{"run_at", "TIMESTAMP WITHOUT TIMEZONE", "NOT NULL DEFAULT CURRENT_TIMESTAMP"}

	switch m.db.Dialect {
	case sqlite3.Dialect:
		ver.datatype = "TEXT"
		state.datatype = "TEXT"
		runat.datatype = "DATETIME"
	case mssql.Dialect, sqlserver.Dialect:
		runat.datatype = "DATETIME"
	case postgresql.Dialect, mysql.Dialect:
	default:
		return fmt.Errorf("%s dialect is not supported by migrator", m.db.Dialect)
	}

	columns := []column{ver, state, runat}
	clmns := []string{}
	for _, c := range columns {
		clmns = append(clmns, fmt.Sprintf("%s %s %s", c.column, c.datatype, c.attributes))
	}
	query := fmt.Sprintf("CREATE TABLE %s (%s)", migrationTable, strings.Join(clmns, ", "))
	_, err := m.db.Exec(query)
	return err
}
