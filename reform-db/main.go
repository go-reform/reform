package main

import (
	"bytes"
	"database/sql"
	"flag"
	"io/ioutil"
	"log"

	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

var (
	fF      = flag.String("f", "", "file to execute")
	driverF = flag.String("db-driver", "", "database driver")
	sourceF = flag.String("db-source", "", "database connection string")
)

func main() {
	log.SetPrefix("reform-db: ")
	flag.Parse()
	log.Print("Internal tool. Do not use it yet.")

	b, err := ioutil.ReadFile(*fF)
	if err != nil {
		log.Fatal(err)
	}

	db, err := sql.Open(*driverF, *sourceF)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	b = bytes.TrimSpace(b)
	if len(b) > 0 {
		_, err = db.Exec(string(b))
		if err != nil {
			log.Fatal(err)
		}
	}
}
