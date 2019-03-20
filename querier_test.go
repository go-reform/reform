package reform_test

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"testing"
	"time"

	mssqlDriver "github.com/denisenkom/go-mssqldb"
	mysqlDriver "github.com/go-sql-driver/mysql"
	"github.com/jackc/pgx/stdlib"
	"github.com/lib/pq"
	sqlite3Driver "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gopkg.in/reform.v1"
	"gopkg.in/reform.v1/dialects/mssql"
	"gopkg.in/reform.v1/dialects/mysql"
	"gopkg.in/reform.v1/dialects/postgresql"
	"gopkg.in/reform.v1/dialects/sqlite3"
	"gopkg.in/reform.v1/dialects/sqlserver"
)

func sleepQuery(t testing.TB, q *reform.Querier, d time.Duration) string {
	switch q.Dialect {
	case postgresql.Dialect:
		return fmt.Sprintf("SELECT pg_sleep(%f)", d.Seconds())
	case mysql.Dialect:
		return fmt.Sprintf("SELECT SLEEP(%f)", d.Seconds())
	case sqlite3.Dialect:
		return fmt.Sprintf("SELECT sleep(%d)", d.Nanoseconds())
	case mssql.Dialect, sqlserver.Dialect:
		sec := int(d.Seconds())
		msec := (d - time.Duration(sec)*time.Second) / time.Millisecond
		return fmt.Sprintf("WAITFOR DELAY '00:00:%02d.%03d'", sec, msec)
	default:
		t.Fatalf("No sleep for %s.", q.Dialect)
		return ""
	}
}

func TestExecContext(t *testing.T) {
	tx, _ := setupTest(t)
	require.NoError(t, tx.Rollback())
	defer tearDownTest(t, tx, nil)

	tx, err := DB.Begin()
	require.NoError(t, err)
	defer tx.Rollback()

	dbDriver := DB.DBInterface().(*sql.DB).Driver()
	const sleep = 200 * time.Millisecond
	const ctxTimeout = 100 * time.Millisecond
	query := sleepQuery(t, tx.Querier, sleep)
	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
	defer cancel()

	start := time.Now()
	_, err = tx.WithContext(ctx).Exec(query)
	dur := time.Since(start)
	switch dbDriver.(type) {
	case *sqlite3Driver.SQLiteDriver:
		assert.NoError(t, err)
		assert.True(t, dur >= sleep, "sqlite3: dur < sleep")
		assert.True(t, dur >= ctxTimeout, "sqlite3: dur < ctxTimeout")
	default:
		assert.Error(t, err)
		assert.True(t, dur < sleep, "dur >= sleep")
		assert.True(t, dur > ctxTimeout, "dur <= ctxTimeout")
	}

	var res int
	err = tx.QueryRow("SELECT 1").Scan(&res)
	switch dbDriver.(type) {
	case *pq.Driver:
		require.Error(t, err)
		assert.Contains(t, err.Error(), `current transaction is aborted, commands ignored until end of transaction block`)
	case *stdlib.Driver:
		require.Error(t, err)
		assert.Contains(t, err.Error(), `current transaction is aborted, commands ignored until end of transaction block`)
	case *mysqlDriver.MySQLDriver:
		assert.Equal(t, driver.ErrBadConn, err)
	case *sqlite3Driver.SQLiteDriver:
		assert.NoError(t, err)
	case *mssqlDriver.Driver:
		assert.Equal(t, driver.ErrBadConn, err)
	default:
		t.Fatalf("QueryRow: unhandled driver %T. err = %s", dbDriver, err)
	}

	err = tx.Rollback()
	switch dbDriver.(type) {
	case *pq.Driver:
		assert.NoError(t, err)
	case *stdlib.Driver:
		assert.NoError(t, err)
	case *mysqlDriver.MySQLDriver:
		assert.Equal(t, mysqlDriver.ErrInvalidConn, err)
	case *sqlite3Driver.SQLiteDriver:
		assert.NoError(t, err)
	case *mssqlDriver.Driver:
		assert.Equal(t, driver.ErrBadConn, err)
	default:
		t.Fatalf("Rollback: unhandled driver %T. err = %s", dbDriver, err)
	}
}
