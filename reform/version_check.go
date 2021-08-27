//go:build !go1.17
// +build !go1.17

package main

import (
	"log"
	"runtime"
)

func init() {
	log.Fatalf("reform requires Go 1.17+, but was compiled with %s.", runtime.Version())
}
