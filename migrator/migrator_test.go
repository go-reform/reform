package migrator

import (
	"os"
	"testing"

	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"gopkg.in/reform.v1/migrator/models"

	"github.com/stretchr/testify/require"
	"gopkg.in/reform.v1"
	"gopkg.in/reform.v1/internal"
)

var (
	DB *reform.DB
)

func TestMain(m *testing.M) {
	DB = internal.ConnectToTestDB()
	os.Exit(m.Run())
}

func TestCreateMigrationTable(t *testing.T) {
	// delete create_migration table if exists
	query := "DROP TABLE IF EXISTS schema_migrations"
	_, err := DB.Exec(query)
	require.NoError(t, err)

	m := Manager{db: DB}
	err = m.createMigrationTable()
	require.NoError(t, err)

	_, err = DB.SelectAllFrom(models.MigrationTable, "")
	require.NoError(t, err)
}
