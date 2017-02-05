package main

import (
	"bytes"
	"io/ioutil"
	"os"

	"gopkg.in/reform.v1"
)

func readFiles(files []string) (queries []string) {
	if len(files) == 0 {
		logger.Debugf("no files are given, reading stdin")
		b, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			logger.Fatalf("failed to read stdin: %s", err)
		}
		b = bytes.TrimSpace(b)
		if len(b) > 0 {
			queries = append(queries, string(b))
		}
	}

	for _, f := range files {
		logger.Debugf("reading file %s", f)
		b, err := ioutil.ReadFile(f)
		if err != nil {
			logger.Fatalf("failed to read file %s: %s", f, err)
		}
		b = bytes.TrimSpace(b)
		if len(b) > 0 {
			queries = append(queries, string(b))
		}
	}

	return
}

func cmdExec(db *reform.DB, files []string) {
	queries := readFiles(files)
	for _, q := range queries {
		_, err := db.Exec(q)
		if err != nil {
			logger.Fatalf("failed to execute %s: %s", q, err)
		}
	}
}
