package main

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"gopkg.in/reform.v1"
)

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

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.Debug)
		fmt.Fprintln(w, strings.Join(columns, "\t"))

		for i, c := range columns {
			columns[i] = strings.Repeat("-", len(c))
		}
		fmt.Fprintln(w, strings.Join(columns, "\t"))

		for rows.Next() {
			line := make([][]byte, len(columns))
			dests := make([]interface{}, len(line))
			for i := range dests {
				dests[i] = &line[i]
			}
			if err = rows.Scan(dests...); err != nil {
				logger.Fatal(err)
			}
			fmt.Fprintf(w, "%s\n", bytes.Join(line, []byte("\t")))
		}
		if err = rows.Err(); err != nil {
			logger.Fatal(err)
		}

		if err = w.Flush(); err != nil {
			logger.Fatal(err)
		}
		if err = rows.Close(); err != nil {
			logger.Fatal(err)
		}
	}
}
