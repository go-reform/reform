package main

import (
	"bytes"
	"io/ioutil"
	"os"

	"gopkg.in/reform.v1"
)

func cmdExec(db *reform.DB, files []string) {
	var b []byte
	if len(files) == 0 {
		var err error
		b, err = ioutil.ReadAll(os.Stdin)
		if err != nil {
			logger.Fatalf("failed to read stdin: %s", err)
		}
	}

	for _, f := range files {
		b1, err := ioutil.ReadFile(f)
		if err != nil {
			logger.Fatalf("failed to read %s: %s", f, err)
		}
		b = append(b, b1...)
	}

	b = bytes.TrimSpace(b)
	if len(b) > 0 {
		_, err := db.Exec(string(b))
		if err != nil {
			logger.Fatalf("failed to execute %s: %s", b, err)
		}
	}
}

func readSQL(path string) ([]byte, error) {
	if path == "" {
		return ioutil.ReadAll(os.Stdin)
	}
	return ioutil.ReadFile(path)
}
