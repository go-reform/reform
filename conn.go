package reform

import (
	"context"
	"database/sql"
	"time"
)

// connInterface is a subset of *sql.Conn used by reform.
// Can be used together with NewConnFromInterface for easier integration with existing code or for passing test doubles.
//
// It may grow and shrink over time to include only needed *sql.Conn methods,
// and is excluded from SemVer compatibility guarantees.
type connInterface interface {
	DBTXContext
	Close() error
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
}

// check interface
var _ connInterface = (*sql.Conn)(nil)

type Conn struct {
	*Querier
	conn connInterface
}

func newConn(ctx context.Context, conn connInterface, dialect Dialect, logger Logger) *Conn {
	return &Conn{
		Querier: newQuerier(ctx, conn, "", dialect, logger),
		conn:    conn,
	}
}

// Begin starts transaction with Querier's context and default options.
func (c *Conn) Begin() (*TX, error) {
	return c.BeginTx(c.ctx, nil)
}

// BeginTx starts transaction with given context and options (can be nil).
func (c *Conn) BeginTx(ctx context.Context, opts *sql.TxOptions) (*TX, error) {
	c.logBefore("BEGIN", nil)
	start := time.Now()
	tx, err := c.conn.BeginTx(ctx, opts)
	c.logAfter("BEGIN", nil, time.Since(start), err)
	if err != nil {
		return nil, err
	}
	return newTX(ctx, tx, c.Dialect, c.Logger), nil
}

func (c *Conn) Close() error {
	return c.conn.Close()
}

// check interfaces
var (
	_ DBTX        = (*Conn)(nil)
	_ DBTXContext = (*Conn)(nil)
)
