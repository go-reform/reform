package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	migrateFlags = flag.NewFlagSet("migrate", flag.ExitOnError)
)

func init() {
	migrateFlags.Usage = func() {
		fmt.Fprintf(os.Stderr, "`migrate` command executes SQL migrations from given files or stdin.\n\n")
		fmt.Fprintf(os.Stderr, "Usage:\n")
		fmt.Fprintf(os.Stderr, "  %s [global flags] migrate [file names]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Global flags:\n")
		flag.PrintDefaults()
		migrateFlags.PrintDefaults()
		fmt.Fprintf(os.Stderr, `
Each migration should be started with a magic version comment "-- +migration VERSION",
where VERSION should match filename prefix. Each migration is executed as a single transaction.
If file names are not given, a migration is read from stdin until EOF, then executed.
`)
	}
}

func cmdMigrate() error {
	return nil
}
