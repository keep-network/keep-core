package main

import (
	"fmt"
	"os"
	"time"

	"github.com/dfinity/go-dfinity-crypto/bls"
	"github.com/keep-network/keep-core/cmd"
	"github.com/urfave/cli"
)

const (
	defaultConfigPath string = "./config.toml"
	defaultGroupSize  int    = 10
	defaultThreshold  int    = 4
)

var (
	cmds       []cli.Command
	configPath string
	// GroupSize ...
	GroupSize int
	// Threshold ...
	Threshold int
	// Version is the semantic version (added at compile time)  See scripts/version.sh
	Version string

	// Revision is the git commit id (added at compile time)
	Revision string

	keepCommands = []cli.Command{
		{
			Name:        "smoke-test",
			Usage:       "smoke-test",
			Description: "Simulate DKG (10 members, threshold 4) and verify group's threshold signature",
			Action:      cmd.SmokeTest,
		},
	}
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
			Name:  "Keep Network",
			Email: "info@keep.network",
		},
	}
	app.Copyright = ""
	app.HelpName = "keep-client"
	app.Usage = "The Keep Client Application"
	app.Commands = keepCommands
	app.Action = func(c *cli.Context) error {
		return nil

	}
	app.Run(os.Args)
}
