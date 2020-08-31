package reform

import (
	"context"
	"database/sql"
)

// connInterface is a subset of *sql.Conn used by reform.
// Can be used together with NewConnFromInterface for easier integration with existing code or for passing test doubles.
//
// It may grow and shrink over time to include only needed *sql.Conn methods,
// and is excluded from SemVer compatibility guarantees.
type connInterface interface {
	DBTXContext
	Close() error
	// TODO Begin, BeginTx?
	// TODO Ping?
	// TODO Prepare?
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

func (c *Conn) Close() error {
	return c.conn.Close()
}

// check interfaces
var (
	_ DBTX        = (*Conn)(nil)
	_ DBTXContext = (*Conn)(nil)
)
