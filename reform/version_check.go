//go:build !go1.26

package main

import (
	"log"
	"runtime"
)

func init() {
	log.Fatalf("reform requires Go 1.26+, but was compiled with %s.", runtime.Version())
}
