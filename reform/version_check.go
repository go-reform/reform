// +build !go1.8

package main

import (
	"log"
	"runtime"
)

func init() {
	log.Fatalf("reform requires Go 1.8+, but was compiled with %s.", runtime.Version())
}
