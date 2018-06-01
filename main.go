package main

import (
	"log"
	"os"

	"github.com/dfinity/go-dfinity-crypto/bls"
	"github.com/keep-network/keep-core/cli"
)

var (
	// Version is the semantic version (added at compile time)  See scripts/version.sh
	Version string

	// Revision is the git commit id (added at compile time)
	Revision string
)

func init() {
	//TODO: Remove Version and Revision when build process auto-populates these values
	Version = "0.0.1"
	Revision = "deadbeef"
}

func main() {

	// Initialize BLS library
	err := bls.Init(bls.CurveSNARK1)
	if err != nil {
		log.Fatal("Failed to initialize BLS.", err)
	}

	cliErr := cli.RunCLI(os.Args, Version, Revision)
	if cliErr != nil {
		log.Println("CLI error encountered:")
		log.Fatal(cliErr)
	}

}
