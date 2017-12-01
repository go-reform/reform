package reform

import (
	"database/sql"
	"time"
)

// TXInterface is a subset of *sql.Tx used by reform.
// Can be used together with NewTXFromInterface for easier integration with existing code or for passing test doubles.
type TXInterface interface {
	DBTX
	Commit() error
	Rollback() error
}

// check interface
var _ TXInterface = (*sql.Tx)(nil)

// TX represents a SQL database transaction.
type TX struct {
	*Querier
	tx TXInterface
}

// NewTX creates new TX object for given SQL database transaction.
// Logger can be nil.
func NewTX(tx *sql.Tx, dialect Dialect, logger Logger) *TX {
	return NewTXFromInterface(tx, dialect, logger)
}

// NewTXFromInterface creates new TX object for given TXInterface.
// Can be used for easier integration with existing code or for passing test doubles.
// Logger can be nil.
func NewTXFromInterface(tx TXInterface, dialect Dialect, logger Logger) *TX {
	return &TX{
		Querier: newQuerier(tx, dialect, logger),
		tx:      tx,
	}
}

// Commit commits the transaction.
func (tx *TX) Commit() error {
	tx.logBefore("COMMIT", nil)
	start := time.Now()
	err := tx.tx.Commit()
	tx.logAfter("COMMIT", nil, time.Since(start), err)
	return err
}

// Rollback aborts the transaction.
func (tx *TX) Rollback() error {
	tx.logBefore("ROLLBACK", nil)
	start := time.Now()
	err := tx.tx.Rollback()
	tx.logAfter("ROLLBACK", nil, time.Since(start), err)
	return err
}

// check interface
var _ DBTX = (*TX)(nil)
