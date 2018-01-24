package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"gopkg.in/reform.v1/migrator"
)

var (
	createMigraitonFlags = flag.NewFlagSet("create-migration", flag.ExitOnError)
	migrationDir         = createMigraitonFlags.String("dir", "migrations", "Directory to store migrations")
)

func init() {
	createMigraitonFlags.Usage = func() {
		fmt.Fprintf(os.Stderr, "`create-migration` prepares a new migration file.\n\n")
		fmt.Fprintf(os.Stderr, "Usage:\n")
		fmt.Fprintf(os.Stderr, "  %s [global flags] create-migration migration-name\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Global flags:\n")
		flag.PrintDefaults()
		createMigraitonFlags.PrintDefaults()
		fmt.Fprintf(os.Stderr, `
It creates a file <2006-01-02-15-04>-migration-name.sql
where <2006-01-02-15-04> is current date and migration-name is given file name.
`)
	}
}

func cmdCreateMigration() error {
	if createMigraitonFlags.NArg() != 1 {
		return fmt.Errorf("Expected one argument for %q, got %d", "create-migration", initFlags.NArg())
	}

	migration := createMigraitonFlags.Arg(0)
	_, err := migrator.Create(*migrationDir, migration, time.Now().UTC())
	return err
}
