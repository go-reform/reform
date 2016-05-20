package reform

import (
	"database/sql"
	"time"
)

// DB represents a connection to SQL database.
type DB struct {
	*Querier
	db *sql.DB
}

// NewDB creates new DB object for given SQL database connection.
func NewDB(db *sql.DB, dialect Dialect, logger Logger) *DB {
	return &DB{
		Querier: newQuerier(db, dialect, logger),
		db:      db,
	}
}

// Begin starts a transaction.
func (db *DB) Begin() (*TX, error) {
	start := time.Now()
	db.logBefore("BEGIN", nil)
	tx, err := db.db.Begin()
	db.logAfter("BEGIN", nil, time.Now().Sub(start), err)
	if err != nil {
		return nil, err
	}
	return NewTX(tx, db.Dialect, db.Logger), nil
}

// InTransaction wraps function execution in transaction, rolling back it in case of error or panic,
// committing otherwise.
func (db *DB) InTransaction(f func(t *TX) error) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	var committed bool
	defer func() {
		if !committed {
			tx.Rollback()
		}
	}()

	err = f(tx)
	if err == nil {
		err = tx.Commit()
	}
	if err == nil {
		committed = true
	}
	return err
}

// check interface
var _ DBTX = new(DB)
