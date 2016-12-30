package reform

import (
	"database/sql"
	"strconv"
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
	tx        TXInterface
	savepoint int
}

// NewTX creates new TX object for given SQL database transaction.
func NewTX(tx *sql.Tx, dialect Dialect, logger Logger) *TX {
	return NewTXFromInterface(tx, dialect, logger)
}

// NewTXFromInterface creates new TX object for given TXInterface.
// Can be used for easier integration with existing code or for passing test doubles.
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

func (tx *TX) Savepoint() error {
	tx.savepoint++
	query := "SAVEPOINT reform_" + strconv.Itoa(tx.savepoint)
	_, err := tx.Exec(query)
	if err != nil {
		tx.savepoint--
	}
	return err
}

func (tx *TX) ReleaseSavepoint() error {
	if tx.savepoint == 0 {
		return ErrNoSavepoint
	}

	query := "RELEASE SAVEPOINT reform_" + strconv.Itoa(tx.savepoint)
	tx.savepoint--
	_, err := tx.Exec(query)
	if err != nil {
		tx.savepoint++
	}
	return err
}

func (tx *TX) RollbackToSavepoint() error {
	if tx.savepoint == 0 {
		return ErrNoSavepoint
	}

	query := "ROLLBACK TO SAVEPOINT reform_" + strconv.Itoa(tx.savepoint)
	tx.savepoint--
	_, err := tx.Exec(query)
	if err != nil {
		tx.savepoint++
	}
	return err
}

// InSavepoint wraps function execution in savepoint, rolling back it in case of error or panic,
// committing (releasing) otherwise.
func (tx *TX) InSavepoint(f func() error) error {
	err := tx.Savepoint()
	if err != nil {
		return err
	}

	var released bool
	defer func() {
		if !released {
			// always return f() or ReleaseSavepoint() error, not possible RollbackSavepoint() error
			_ = tx.RollbackToSavepoint()
		}
	}()

	err = f()
	if err == nil {
		err = tx.ReleaseSavepoint()
	}
	if err == nil {
		released = true
	}
	return err
}

// check interface
var _ DBTX = (*TX)(nil)
