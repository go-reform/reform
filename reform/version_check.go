//go:build !go1.15
// +build !go1.15

package main

import (
	"log"
	"runtime"
)

func init() {
	log.Fatalf("reform requires Go 1.15+, but was compiled with %s.", runtime.Version())
}
