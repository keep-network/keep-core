package main

import (
	"log"
	"os"

	"fmt"
	"path"
	"time"

	"github.com/dfinity/go-dfinity-crypto/bls"
	"github.com/keep-network/keep-core/cmd"
	"github.com/urfave/cli"
)

var (
	// Version is the CLI semver version
	Version string

	// Revision is the git commit (revision) hash
	Revision string

	// GroupSize indicates the number of members in this relay group
	GroupSize int

	// Threshold indicates the threshold number of members required to perform signature verification
	Threshold int

	configPath string
)

const (
	defaultGroupSize int = 10
	defaultThreshold int = 4
)

func init() {
	// Version and Revision should be set by go linker.
	if Version == "" {
		Version = "unknown"
	}
	if Revision == "" {
		Revision = "unknown"
	}
}

func main() {
	err := NewApp(Version, Revision).Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

// NewApp creates a new keep cli application with the respective commands and metainfo.
func NewApp(version, revision string) *cli.App {
	Version = version
	Revision = revision

	app := cli.NewApp()
	app.Name = path.Base(os.Args[0])
	app.Usage = "CLI for The Keep Network"
	app.Description = "Command line interface (CLI) for running a Keep provider"
	app.Copyright = "" //TODO: Insert copyright info later
	app.Compiled = time.Now()
	app.Authors = []cli.Author{
		{
			Name:  "Keep Network",
			Email: "info@keep.network",
		},
	}
	app.Version = fmt.Sprintf("%s (revision %s)", version, revision)
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "config,c",
			Value:       "",
			Destination: &configPath,
			EnvVar:      "KEEP_CONFIG_PATH",
			Usage:       "optionally, specify the environment variable",
		},
	}
	app.Before = func(c *cli.Context) error {
		err := bls.Init(bls.CurveSNARK1)
		if err != nil {
			log.Fatal("Failed to initialize BLS.", err)
		}
		return nil
	}
	app.Commands = []cli.Command{
		{
			Name:      "get-info",
			Usage:     "prints keep client information",
			ArgsUsage: " ", // no args
			Category:  "keep client information",
			Action: func(c *cli.Context) error {
				getInfo(c)
				return nil
			},
		},
		{
			Name:        "smoke-test",
			Usage:       "Simulates DKG and signature verification",
			Description: "Simulate Distributed Key Generation (DKG) and verify group's threshold signature",
			Action:      cmd.SmokeTest,
			Flags:       cmd.SmokeTestFlags,
		},
	}

	cli.AppHelpTemplate = fmt.Sprintf(`%s
ENVIRONMENT VARIABLES:
   KEEP_ETHEREUM_PASSWORD    keep client password

`, cli.AppHelpTemplate)

	return app
}

//-------------------------------------------------------------------------------
// Helpers
//-------------------------------------------------------------------------------

func getInfo(c *cli.Context) {
	fmt.Printf("Keep client: %s\n\n"+
		"Description: %s\n"+
		"Version:     %s\n"+
		"Revision:    %s\n"+
		"Config Path: %s\n",
		c.App.Name, c.App.Description, Version, Revision, c.GlobalString("config"))
}
