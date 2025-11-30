//go:build !go1.24

package main

import (
	"log"
	"runtime"
)

func init() {
	log.Fatalf("reform requires Go 1.24+, but was compiled with %s.", runtime.Version())
}
