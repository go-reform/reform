package main

import (
	"bytes"
	"database/sql"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"

	"gopkg.in/reform.v1"
)

var (
	fF      = flag.String("f", "", "file to execute")
	driverF = flag.String("db-driver", "", "database driver")
	sourceF = flag.String("db-source", "", "database connection string")
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "reform-db. %s.\n\n", reform.Version)
		flag.PrintDefaults()
	}
	flag.Parse()

	log.SetPrefix("reform-db: ")
	log.Print("Internal tool. Do not use it yet.")

	b, err := readSQL(*fF)
	if err != nil {
		log.Fatalf("failed to read %q: %s", *fF, err)
	}

	db, err := sql.Open(*driverF, *sourceF)
	if err != nil {
		log.Fatalf("failed to connect to %s %q: %s", *driverF, *sourceF, err)
	}
	defer db.Close()

	// Use single connection so various session-related variables work.
	// For example: "PRAGMA foreign_keys" for SQLite3, "SET IDENTITY_INSERT" for MS SQL, etc.
	db.SetMaxIdleConns(1)
	db.SetMaxOpenConns(1)
	db.SetConnMaxLifetime(0)

	err = db.Ping()
	if err != nil {
		log.Fatalf("failed to ping database: %s", err)
	}

	b = bytes.TrimSpace(b)
	if len(b) > 0 {
		q := string(b)
		_, err := db.Exec(q)
		if err != nil {
			log.Fatalf("failed to execute %s: %s", q, err)
		}
	}
}

func readSQL(path string) ([]byte, error) {
	if path == "" {
		return ioutil.ReadAll(os.Stdin)
	}
	return ioutil.ReadFile(path)
}
