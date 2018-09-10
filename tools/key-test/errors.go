package main

import (
	"fmt"
	"io"
	"os"
	"runtime"

	"github.com/pschlump/godebug"
)

var errorOutput = io.MultiWriter(os.Stdout, os.Stderr)

func init() {
	if runtime.GOOS == "windows" {
		// The SameFile check below doesn't work on Windows.
		// stdout is unlikely to get redirected though, so just print there.
		errorOutput = os.Stdout
	} else {
		outf, _ := os.Stdout.Stat()
		errf, _ := os.Stderr.Stat()
		if outf != nil && errf != nil && os.SameFile(outf, errf) {
			errorOutput = os.Stderr
		}
	}
}

// Fatalf formats a message to standard error and exits the program.
// On Linux/Mac the message is also printed to standard output if standard error
// is redirected to a different file.
func Fatalf(rc int, format string, args ...interface{}) {
	fmt.Fprintf(errorOutput, "Fatal: "+format+"\n", args...)
	fmt.Fprintf(errorOutput, " From: %s\n", godebug.LF(2))
	os.Exit(rc)
}
