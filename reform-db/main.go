package main

import (
	"database/sql"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"

	"gopkg.in/reform.v1"
	"gopkg.in/reform.v1/internal"
)

var (
	logger *internal.Logger

	debugF  = flag.Bool("debug", false, "Enable debug logging")
	driverF = flag.String("db-driver", "", "database driver")
	sourceF = flag.String("db-source", "", "database connection string")
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "reform-db. %s.\n\n", reform.Version)
		fmt.Fprintf(os.Stderr, "Usage:\n\n")
		fmt.Fprintf(os.Stderr, "  %s [flags] [command] [arguments]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Commands:\n")
		fmt.Fprintf(os.Stderr, "  exec\n")
		fmt.Fprintf(os.Stderr, "  query\n")
		fmt.Fprintf(os.Stderr, "  init\n\n")
		fmt.Fprintf(os.Stderr, "Flags:\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	logger = internal.NewLogger("reform-db: ", *debugF)
	logger.Printf("Internal tool. Do not use it yet.")

	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(1)
	}

	sqlDB, err := sql.Open(*driverF, *sourceF)
	if err != nil {
		logger.Fatalf("failed to connect to %s %q: %s", *driverF, *sourceF, err)
	}
	defer sqlDB.Close()

	// Use single connection so various session-related variables work.
	// For example: "PRAGMA foreign_keys" for SQLite3, "SET IDENTITY_INSERT" for MS SQL, etc.
	sqlDB.SetMaxIdleConns(1)
	sqlDB.SetMaxOpenConns(1)
	sqlDB.SetConnMaxLifetime(0)

	err = sqlDB.Ping()
	if err != nil {
		logger.Fatalf("failed to ping database: %s", err)
	}

	dialect := internal.DialectForDriver(*driverF)
	db := reform.NewDB(sqlDB, dialect, reform.NewPrintfLogger(logger.Debugf))

	switch flag.Arg(0) {
	case "exec":
		cmdExec(db, flag.Args()[1:])

	case "query":
		cmdQuery(db, flag.Args()[1:])

	case "init":
		if flag.NArg() > 1 {
			logger.Fatalf("Expected zero or one argument for %q, got %d", "init", flag.NArg())
		}

		dir := flag.Arg(1)
		if dir == "" {
			if dir, err = os.Getwd(); err != nil {
				logger.Fatalf("%s", err)
			}
		}
		if dir, err = filepath.Abs(dir); err != nil {
			logger.Fatalf("%s", err)
		}
		fi, err := os.Stat(dir)
		if err != nil {
			logger.Fatalf("%s", err)
		}
		if !fi.IsDir() {
			logger.Fatalf("%q should be existing directory", dir)
		}

		cmdInit(db, dir)

	default:
		logger.Fatalf("Unexpected command %q", flag.Arg(0))
	}
}
