package main

import (
	"github.com/dfinity/go-dfinity-crypto/bls"
	"github.com/keep-network/keep-core/cmd/keep"
)

var (
	// Version is the semantic version (added at compile time)  See scripts/version.sh
	Version string

	// Revision is the git commit id (added at compile time)
	Revision string
)

func main() {
	bls.Init(bls.CurveSNARK1)

	keep.ClientApp(Version, Revision)
}
