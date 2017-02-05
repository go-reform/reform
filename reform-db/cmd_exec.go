package main

import (
	"bytes"
	"io/ioutil"
	"os"

	"gopkg.in/reform.v1"
)

func cmdExec(db *reform.DB, files []string) {
	var query []byte
	if len(files) == 0 {
		b, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			logger.Fatalf("failed to read stdin: %s", err)
		}
		query = append(query, b...)
	}

	for _, f := range files {
		b, err := ioutil.ReadFile(f)
		if err != nil {
			logger.Fatalf("failed to read %s: %s", f, err)
		}
		query = append(query, b...)
	}

	query = bytes.TrimSpace(query)
	if len(query) > 0 {
		_, err := db.Exec(string(query))
		if err != nil {
			logger.Fatalf("failed to execute %s: %s", query, err)
		}
	}
}
