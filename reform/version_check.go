// +build !go1.12

package main

import (
	"log"
	"runtime"
)

func init() {
	log.Fatalf("reform requires Go 1.12+, but was compiled with %s.", runtime.Version())
}
