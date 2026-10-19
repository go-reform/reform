package reform_test

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"testing"
	"time"

	mysqlDriver "github.com/go-sql-driver/mysql"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/lib/pq"
	"github.com/lib/pq/pqerror"
	sqlite3Driver "github.com/mattn/go-sqlite3"
	mssqlDriver "github.com/microsoft/go-mssqldb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gopkg.in/reform.v1"
	"gopkg.in/reform.v1/dialects/mssql" //nolint:staticcheck
	"gopkg.in/reform.v1/dialects/mysql"
	"gopkg.in/reform.v1/dialects/postgresql"
	"gopkg.in/reform.v1/dialects/sqlite3"
	"gopkg.in/reform.v1/dialects/sqlserver"
)

type ctxKey string

func sleepQuery(t testing.TB, q *reform.Querier, d time.Duration) string {
	switch q.Dialect {
	case postgresql.Dialect:
		return fmt.Sprintf("SELECT pg_sleep(%f)", d.Seconds())
	case mysql.Dialect:
		return fmt.Sprintf("SELECT SLEEP(%f)", d.Seconds())
	case sqlite3.Dialect:
		return fmt.Sprintf("SELECT sleep(%d)", d.Nanoseconds())
	case mssql.Dialect: //nolint:staticcheck
		fallthrough
	case sqlserver.Dialect:
		sec := int(d.Seconds())
		msec := (d - time.Duration(sec)*time.Second) / time.Millisecond
		return fmt.Sprintf("WAITFOR DELAY '00:00:%02d.%03d'", sec, msec)
	default:
		t.Fatalf("No sleep for %s.", q.Dialect)
		return ""
	}
}

func TestExecWithContext(t *testing.T) {
	db, tx := setupTX(t)
	defer teardown(t, db)

	assert.Equal(t, context.Background(), db.Context())
	assert.Equal(t, context.Background(), tx.Context())

	dbDriver := db.DBInterface().(*sql.DB).Driver()
	const sleep = 200 * time.Millisecond
	const ctxTimeout = 100 * time.Millisecond
	query := sleepQuery(t, tx.Querier, sleep)
	ctx, cancel := context.WithTimeout(context.WithValue(context.Background(), ctxKey("k"), "exec"), ctxTimeout)
	defer cancel()

	q := tx.WithContext(ctx)
	assert.Equal(t, ctx, q.Context())
	start := time.Now()
	_, err := q.Exec(query)
	dur := time.Since(start)

	switch dbDriver.(type) {
	case *sqlite3Driver.SQLiteDriver:
		// SQLite cannot interrupt this Go function; the statement completes successfully.
		assert.NoError(t, err)
		assert.ErrorIs(t, ctx.Err(), context.DeadlineExceeded)
		assert.True(t, dur >= sleep, "sqlite3: failed comparison: dur >= sleep")
		assert.True(t, dur >= ctxTimeout, "sqlite3: failed comparison: dur >= ctxTimeout")

	default:
		assert.Error(t, err)
		assert.True(t, dur < sleep, "failed comparison: dur < sleep")
		assert.True(t, dur > ctxTimeout, "failed comparison: dur > ctxTimeout")

		// check specific error type
		switch dbDriver.(type) {
		case *pq.Driver:
			require.EqualError(t, err, "pq: canceling statement due to user request (57014)")
			pgErr := err.(*pq.Error)
			assert.Equal(t, "ERROR", pgErr.Severity)
			assert.Equal(t, pqerror.Code("57014"), pgErr.Code)
			assert.Equal(t, "ProcessInterrupts", pgErr.Routine)
		case *stdlib.Driver:
			assert.ErrorIs(t, err, context.DeadlineExceeded)
		case *mysqlDriver.MySQLDriver:
			assert.Equal(t, context.DeadlineExceeded, err)
		case *mssqlDriver.Driver:
			assert.Equal(t, context.DeadlineExceeded, err)
		default:
			t.Fatalf("q.Exec: unhandled driver %T. err = %s", dbDriver, err)
		}
	}

	// context should not be modified
	assert.Equal(t, context.Background(), db.Context())
	assert.Equal(t, context.Background(), tx.Context())

	// check q with expired timeout
	var res int
	err = q.QueryRow("SELECT 1").Scan(&res)
	assert.Equal(t, context.DeadlineExceeded, err)
	require.Equal(t, 0, res)

	// check the same tx without timeout
	err = tx.QueryRow("SELECT 1").Scan(&res)
	switch dbDriver.(type) {
	case *pq.Driver, *stdlib.Driver:
		// Both PostgreSQL drivers discard the connection after cancellation.
		assert.Equal(t, driver.ErrBadConn, err)
		assert.Equal(t, 0, res)

	case *mysqlDriver.MySQLDriver:
		assert.Equal(t, driver.ErrBadConn, err)
		assert.Equal(t, 0, res)

	case *sqlite3Driver.SQLiteDriver:
		assert.NoError(t, err)
		assert.Equal(t, 1, res)

	case *mssqlDriver.Driver:
		assert.Equal(t, driver.ErrBadConn, err)
		assert.Equal(t, 0, res)

	default:
		t.Fatalf("tx.QueryRow: unhandled driver %T. err = %s", dbDriver, err)
	}

	err = tx.Rollback()
	switch dbDriver.(type) {
	case *pq.Driver:
		assert.Equal(t, driver.ErrBadConn, err)
	case *stdlib.Driver:
		assert.EqualError(t, err, "conn closed")
	case *mysqlDriver.MySQLDriver:
		assert.ErrorIs(t, err, context.DeadlineExceeded)
	case *sqlite3Driver.SQLiteDriver:
		assert.NoError(t, err)
	case *mssqlDriver.Driver:
		assert.Equal(t, driver.ErrBadConn, err)
	default:
		t.Fatalf("tx.Rollback: unhandled driver %T. err = %s", dbDriver, err)
	}
}
