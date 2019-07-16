package reform_test

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/AlekSi/pointer"
	"github.com/brianvoe/gofakeit"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/stdlib"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gopkg.in/reform.v1"
	"gopkg.in/reform.v1/dialects/postgresql"
	. "gopkg.in/reform.v1/internal/test/models"
)

func TestBeginCommit(t *testing.T) {
	db := setupDB(t)
	defer teardown(t, db)

	person := &Person{ID: 42, Email: pointer.ToString(gofakeit.Email())}

	tx, err := db.Begin()
	require.NoError(t, err)
	assert.NoError(t, insertPersonWithID(t, tx.Querier, person))
	assert.NoError(t, tx.Commit())
	assert.Equal(t, tx.Commit(), reform.ErrTxDone)
	assert.Equal(t, tx.Rollback(), reform.ErrTxDone)
	assert.NoError(t, db.Reload(person))
	assert.NoError(t, db.Delete(person))
}

func TestBeginRollback(t *testing.T) {
	db := setupDB(t)
	defer teardown(t, db)

	person := &Person{ID: 42, Email: pointer.ToString(gofakeit.Email())}

	tx, err := db.Begin()
	require.NoError(t, err)
	assert.NoError(t, insertPersonWithID(t, tx.Querier, person))
	assert.NoError(t, tx.Rollback())
	assert.Equal(t, tx.Commit(), reform.ErrTxDone)
	assert.Equal(t, tx.Rollback(), reform.ErrTxDone)
	assert.Equal(t, db.Reload(person), reform.ErrNoRows)
}

// This behavior is checked for documentation purposes only. reform does not rely on it.
func TestErrorInTransaction(t *testing.T) {
	if DB.Dialect == postgresql.Dialect {
		t.Skip(DB.Dialect.String() + " works differently, see TestAbortedTransaction")
	}

	db := setupDB(t)
	defer teardown(t, db)

	person1 := &Person{ID: 42, Email: pointer.ToString(gofakeit.Email())}
	person2 := &Person{ID: 43, Email: pointer.ToString(gofakeit.Email())}

	// commit works
	tx, err := db.Begin()
	require.NoError(t, err)
	assert.NoError(t, insertPersonWithID(t, tx.Querier, person1))
	assert.Error(t, insertPersonWithID(t, tx.Querier, person1))   // duplicate PK
	assert.NoError(t, insertPersonWithID(t, tx.Querier, person2)) // INSERT works
	assert.NoError(t, tx.Commit())
	assert.Equal(t, tx.Commit(), reform.ErrTxDone)
	assert.Equal(t, tx.Rollback(), reform.ErrTxDone)
	assert.NoError(t, db.Reload(person1))
	assert.NoError(t, db.Reload(person2))
	assert.NoError(t, db.Delete(person1))
	assert.NoError(t, db.Delete(person2))

	// rollback works
	tx, err = db.Begin()
	require.NoError(t, err)
	assert.NoError(t, insertPersonWithID(t, tx.Querier, person1))
	assert.Error(t, insertPersonWithID(t, tx.Querier, person1))   // duplicate PK
	assert.NoError(t, insertPersonWithID(t, tx.Querier, person2)) // INSERT works
	assert.NoError(t, tx.Rollback())
	assert.Equal(t, tx.Commit(), reform.ErrTxDone)
	assert.Equal(t, tx.Rollback(), reform.ErrTxDone)
	assert.EqualError(t, db.Reload(person1), reform.ErrNoRows.Error())
	assert.EqualError(t, db.Reload(person2), reform.ErrNoRows.Error())
}

// This behavior is checked for documentation purposes only. reform does not rely on it.
// http://postgresql.nabble.com/Current-transaction-is-aborted-commands-ignored-until-end-of-transaction-block-td5109252.html
func TestAbortedTransaction(t *testing.T) {
	if DB.Dialect != postgresql.Dialect {
		t.Skip(DB.Dialect.String() + " works differently, see TestErrorInTransaction")
	}

	db := setupDB(t)
	defer teardown(t, db)

	person1 := &Person{ID: 42, Email: pointer.ToString(gofakeit.Email())}
	person2 := &Person{ID: 43, Email: pointer.ToString(gofakeit.Email())}

	// commit fails
	tx, err := db.Begin()
	require.NoError(t, err)
	assert.NoError(t, insertPersonWithID(t, tx.Querier, person1))
	assert.Contains(t, insertPersonWithID(t, tx.Querier, person1).Error(), `duplicate key value violates unique constraint "people_pkey"`)
	assert.Contains(t, insertPersonWithID(t, tx.Querier, person2).Error(), `current transaction is aborted, commands ignored until end of transaction block`)
	err = tx.Commit()
	require.Error(t, err)

	switch db.DBInterface().(*sql.DB).Driver().(type) {
	case *pq.Driver:
		assert.Equal(t, pq.ErrInFailedTransaction, err)
	case *stdlib.Driver:
		assert.Equal(t, pgx.ErrTxCommitRollback, err)
	default:
		t.Fatalf("unexpected driver, error %v", err)
	}
	assert.Equal(t, tx.Rollback(), reform.ErrTxDone)
	assert.EqualError(t, db.Reload(person1), reform.ErrNoRows.Error())
	assert.EqualError(t, db.Reload(person2), reform.ErrNoRows.Error())

	// rollback works
	tx, err = db.Begin()
	require.NoError(t, err)
	assert.NoError(t, insertPersonWithID(t, tx.Querier, person1))
	assert.Contains(t, insertPersonWithID(t, tx.Querier, person1).Error(), `duplicate key value violates unique constraint "people_pkey"`)
	assert.Contains(t, insertPersonWithID(t, tx.Querier, person2).Error(), `current transaction is aborted, commands ignored until end of transaction block`)
	assert.NoError(t, tx.Rollback())
	assert.Equal(t, tx.Commit(), reform.ErrTxDone)
	assert.Equal(t, tx.Rollback(), reform.ErrTxDone)
	assert.EqualError(t, db.Reload(person1), reform.ErrNoRows.Error())
	assert.EqualError(t, db.Reload(person2), reform.ErrNoRows.Error())
}

func TestInTransaction(t *testing.T) {
	db := setupDB(t)
	defer teardown(t, db)

	person := &Person{ID: 42, Email: pointer.ToString(gofakeit.Email())}

	// error in closure
	err := db.InTransaction(func(tx *reform.TX) error {
		assert.NoError(t, insertPersonWithID(t, tx.Querier, person))
		return errors.New("epic error")
	})
	assert.EqualError(t, err, "epic error")
	assert.Equal(t, db.Reload(person), reform.ErrNoRows)

	// panic in closure
	assert.Panics(t, func() {
		err = db.InTransaction(func(tx *reform.TX) error {
			assert.NoError(t, insertPersonWithID(t, tx.Querier, person))
			panic("epic panic!")
		})
	})
	assert.Equal(t, db.Reload(person), reform.ErrNoRows)

	// duplicate PK in closure
	err = db.InTransaction(func(tx *reform.TX) error {
		assert.NoError(t, insertPersonWithID(t, tx.Querier, person))
		err = insertPersonWithID(t, tx.Querier, person)
		assert.Error(t, err)
		return err
	})
	assert.Error(t, err)
	assert.Equal(t, db.Reload(person), reform.ErrNoRows)

	// no error
	err = db.InTransaction(func(tx *reform.TX) error {
		assert.NoError(t, insertPersonWithID(t, tx.Querier, person))
		return nil
	})
	assert.NoError(t, err)
	assert.NoError(t, db.Reload(person))
	assert.NoError(t, db.Delete(person))
}
