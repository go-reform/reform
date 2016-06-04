package reform

import (
	"database/sql"
	"time"
)

// Querier performs queries and commands.
type Querier struct {
	dbtx DBTX
	Dialect
	Logger Logger
}

func newQuerier(dbtx DBTX, dialect Dialect, logger Logger) *Querier {
	return &Querier{
		dbtx:    dbtx,
		Dialect: dialect,
		Logger:  logger,
	}
}

func (q *Querier) logBefore(query string, args []interface{}) {
	if q.Logger != nil {
		q.Logger.Before(query, args)
	}
}

func (q *Querier) logAfter(query string, args []interface{}, d time.Duration, err error) {
	if q.Logger != nil {
		q.Logger.After(query, args, d, err)
	}
}

// QualifiedView returns quoted qualified view name.
func (q *Querier) QualifiedView(view View) string {
	v := q.QuoteIdentifier(view.Name())
	if view.Schema() != "" {
		v = q.QuoteIdentifier(view.Schema()) + "." + v
	}
	return v
}

// QualifiedColumns returns a slice of quoted qualified column names for given view.
func (q *Querier) QualifiedColumns(view View) []string {
	v := q.QualifiedView(view)
	res := view.Columns()
	for i := 0; i < len(res); i++ {
		res[i] = v + "." + q.QuoteIdentifier(res[i])
	}
	return res
}

// Exec executes a query without returning any rows.
// The args are for any placeholder parameters in the query.
func (q *Querier) Exec(query string, args ...interface{}) (sql.Result, error) {
	start := time.Now()
	q.logBefore(query, args)
	res, err := q.dbtx.Exec(query, args...)
	q.logAfter(query, args, time.Now().Sub(start), err)
	return res, err
}

// Query executes a query that returns rows, typically a SELECT.
// The args are for any placeholder parameters in the query.
func (q *Querier) Query(query string, args ...interface{}) (*sql.Rows, error) {
	start := time.Now()
	q.logBefore(query, args)
	rows, err := q.dbtx.Query(query, args...)
	q.logAfter(query, args, time.Now().Sub(start), err)
	return rows, err
}

// QueryRow executes a query that is expected to return at most one row.
// QueryRow always returns a non-nil value. Errors are deferred until Row's Scan method is called.
func (q *Querier) QueryRow(query string, args ...interface{}) *sql.Row {
	start := time.Now()
	q.logBefore(query, args)
	row := q.dbtx.QueryRow(query, args...)
	q.logAfter(query, args, time.Now().Sub(start), nil)
	return row
}

// check interface
var _ DBTX = new(Querier)
