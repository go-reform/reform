package internal

import (
	"database/sql"
	"log"
	"os"
	"time"

	"gopkg.in/reform.v1"
	"gopkg.in/reform.v1/dialects"
	"gopkg.in/reform.v1/dialects/mssql"
	"gopkg.in/reform.v1/dialects/mysql"
	"gopkg.in/reform.v1/dialects/postgresql"
	"gopkg.in/reform.v1/dialects/sqlite3"
	"gopkg.in/reform.v1/dialects/sqlserver"
)

// ConnectToTestDB returns open and prepared connection to test DB.
func ConnectToTestDB() *reform.DB {
	driver := os.Getenv("REFORM_DRIVER")
	source := os.Getenv("REFORM_TEST_SOURCE")
	log.Printf("driver = %q, source = %q", driver, source)
	if driver == "" || source == "" {
		log.Fatal("no driver or source, set REFORM_DRIVER and REFORM_TEST_SOURCE")
	}

	db, err := sql.Open(driver, source)
	if err != nil {
		log.Fatal(err)
	}

	// Use single connection so various session-related variables work.
	// For example: "PRAGMA foreign_keys" for SQLite3, "SET IDENTITY_INSERT" for MS SQL, etc.
	db.SetMaxIdleConns(1)
	db.SetMaxOpenConns(1)
	db.SetConnMaxLifetime(0)

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

	// print useful information for debugging
	now := time.Now()
	log.Printf("time.Now()       = %s", now)
	log.Printf("time.Now().UTC() = %s", now.UTC())

	// select dialect for driver
	dialect := dialects.ForDriver(driver)
	switch dialect {
	case postgresql.Dialect:
		var version, tz string
		if err = db.QueryRow("SHOW server_version").Scan(&version); err != nil {
			log.Fatal(err)
		}
		if err = db.QueryRow("SHOW TimeZone").Scan(&tz); err != nil {
			log.Fatal(err)
		}
		log.Printf("PostgreSQL version  = %q", version)
		log.Printf("PostgreSQL TimeZone = %q", tz)

	case mysql.Dialect:
		q := "SELECT @@version, @@sql_mode, @@autocommit, @@time_zone"
		var version, mode, autocommit, tz string
		if err = db.QueryRow(q).Scan(&version, &mode, &autocommit, &tz); err != nil {
			log.Fatal(err)
		}
		log.Printf("MySQL version    = %q", version)
		log.Printf("MySQL sql_mode   = %q", mode)
		log.Printf("MySQL autocommit = %q", autocommit)
		log.Printf("MySQL time_zone  = %q", tz)

	case sqlite3.Dialect:
		var version, source string
		if err = db.QueryRow("SELECT sqlite_version(), sqlite_source_id()").Scan(&version, &source); err != nil {
			log.Fatal(err)
		}
		log.Printf("SQLite3 version = %q", version)
		log.Printf("SQLite3 source  = %q", source)

		if _, err = db.Exec("PRAGMA foreign_keys = ON"); err != nil {
			log.Fatal(err)
		}

	case mssql.Dialect, sqlserver.Dialect:
		var version string
		var options uint16
		if err = db.QueryRow("SELECT @@VERSION, @@OPTIONS").Scan(&version, &options); err != nil {
			log.Fatal(err)
		}
		xact := "ON"
		if options&0x4000 == 0 {
			xact = "OFF"
		}
		log.Printf("MS SQL version = %s", version)
		log.Printf("MS SQL OPTIONS = %#4x (XACT_ABORT %s)", options, xact)

	default:
		log.Fatalf("reform: no dialect for driver %s", driver)
	}

	return reform.NewDB(db, dialect, nil)
}
