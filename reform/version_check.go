// +build !go1.10

package main

import (
	"log"
	"runtime"
)

func init() {
	log.Fatalf("reform requires Go 1.10+, but was compiled with %s.", runtime.Version())
}
