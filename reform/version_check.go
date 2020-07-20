// +build !go1.13

package main

import (
	"log"
	"runtime"
)

func init() {
	log.Fatalf("reform requires Go 1.13+, but was compiled with %s.", runtime.Version())
}
