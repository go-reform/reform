// +build !go1.7

package main

import (
	"log"
	"runtime"
)

func init() {
	log.Fatalf("reform requires Go 1.7+, but was compiled with %s.", runtime.Version())
}
