package util

import "io"

// CloseReadOnlyFile is typically called by defer; notice that the error is intentionally ignored.
func CloseReadOnlyFile(i interface{}) {
	if f, ok := i.(io.Closer); ok {
		// ignore close of read only file
		_ = f.Close()
	}
}
