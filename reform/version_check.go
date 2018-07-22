// +build !go1.9

package main

import (
	"log"
	"runtime"
)

func init() {
	log.Fatalf("reform requires Go 1.9+, but was compiled with %s.", runtime.Version())
}
