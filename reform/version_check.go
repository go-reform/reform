// +build !go1.6

package main

import (
	"log"
	"runtime"
)

func init() {
	log.Fatalf("reform requires Go 1.6+, but was compiled with %s.", runtime.Version())
}
