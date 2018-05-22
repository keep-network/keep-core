package main

import (
	"fmt"
	"os"
	"time"

	"github.com/dfinity/go-dfinity-crypto/bls"
	"github.com/keep-network/keep-core/cmd/keep"
	"github.com/urfave/cli"
)

var (
	// Version is the semantic version (added at compile time)  See scripts/version.sh
	Version string

	// Revision is the git commit id (added at compile time)
	Revision string
)

//TODO: Remove init when build process is ready to populate Version and Revision
func init() {
	Version = "0.0.1"
	Revision = "deadbeef"
}

func main() {
	bls.Init(bls.CurveSNARK1)

	app := cli.NewApp()
	app.Name = "keep-client"
	app.Version = fmt.Sprintf("%s (revision %s)", Version, Revision)
	app.Compiled = time.Now()
	app.Authors = []cli.Author{
		cli.Author{
			Name:  "Keep Authors",
			Email: "noreply@example.com",
		},
	}
	app.Copyright = ""
	app.HelpName = "keep-client"
	app.Usage = "The Keep Client Application"
	app.Commands = keep.KeepCommands
	app.Action = func(c *cli.Context) error {
		return nil

	}
	app.Run(os.Args)
}
