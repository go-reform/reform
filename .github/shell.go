package main

// Custom shell for GNU Make to measure command execution time.

import (
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

func main() {
	log.SetFlags(0)

	cmd := exec.Command("/bin/bash", os.Args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	start := time.Now()
	err := cmd.Run()
	d := time.Since(start)
	if d > time.Second {
		log.Printf("===> %s: %s", strings.Join(cmd.Args, " "), d)
	}
	if err != nil {
		log.Fatal(err)
	}
}
