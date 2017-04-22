package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"gopkg.in/reform.v1"
)

var (
	queryFlags = flag.NewFlagSet("query", flag.ExitOnError)
)

func init() {
	queryFlags.Usage = func() {
		fmt.Fprintf(os.Stderr, "`query` command executes SQL queries from given files or stdin, and returns results.\n\n")
		fmt.Fprintf(os.Stderr, "Usage:\n")
		fmt.Fprintf(os.Stderr, "  %s [global flags] query [file names]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Global flags:\n")
		flag.PrintDefaults()
		queryFlags.PrintDefaults()
		fmt.Fprintf(os.Stderr, `
Each file's content is executed as a single query. If it contains multiple
statements, make sure SQL driver supports them. If file names are not given,
a query is read from stdin until EOF, then executed.
`)
	}
}

// cmdQuery implements query command.
func cmdQuery(db *reform.DB, files []string) {
	queries := readFiles(files)
	for _, q := range queries {
		rows, err := db.Query(q)
		if err != nil {
			logger.Fatalf("failed to query %s: %s", q, err)
		}
		columns, err := rows.Columns()
		if err != nil {
			logger.Fatalf("failed to get columns for %s: %s", q, err)
		}
		logger.Debugf("result columns: %v", columns)

		// write table header
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.Debug)
		if _, err = fmt.Fprintln(w, strings.Join(columns, "\t")); err != nil {
			logger.Fatalf("%s", err)
		}
		for i, c := range columns {
			columns[i] = strings.Repeat("-", len(c))
		}
		if _, err = fmt.Fprintln(w, strings.Join(columns, "\t")); err != nil {
			logger.Fatalf("%s", err)
		}

		// read all rows, scan each field to []byte
		for rows.Next() {
			line := make([][]byte, len(columns))
			dests := make([]interface{}, len(line))
			for i := range dests {
				dests[i] = &line[i]
			}
			if err = rows.Scan(dests...); err != nil {
				logger.Fatalf("%s", err)
			}
			fmt.Fprintf(w, "%s\n", bytes.Join(line, []byte("\t")))
		}
		if err = rows.Err(); err != nil {
			logger.Fatalf("%s", err)
		}

		if err = w.Flush(); err != nil {
			logger.Fatalf("%s", err)
		}
		if err = rows.Close(); err != nil {
			logger.Fatalf("%s", err)
		}
	}
}
