package main

// This is to make the main and the CLI arguments testable

import (
	"github.com/keep-network/keep-core/go/interface/cli"
)

func main() {
	cli := cli.CLI{}
	cli.Run()
}
