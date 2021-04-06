// +build !go1.14

package main

import (
	"log"
	"runtime"
)

func init() {
	log.Fatalf("reform requires Go 1.14+, but was compiled with %s.", runtime.Version())
}
