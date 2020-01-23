// +build !go1.11

package main

import (
	"log"
	"runtime"
)

func init() {
	log.Fatalf("reform requires Go 1.11+, but was compiled with %s.", runtime.Version())
}
