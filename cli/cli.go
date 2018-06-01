package cli

import (
	"fmt"
	"os"
	"path"
	"time"

	"github.com/keep-network/keep-core/cli/runtime"
	"github.com/urfave/cli"
)

var (
	app                    *cli.App
	configPath             string
	extraArg               string
	blsInitialized         = false
	configWarningDisplayed = false
)

func init() {
	app = cli.NewApp()
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
	app.Commands = commands
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "config,c",
			Value:       "",
			Destination: &configPath,
			EnvVar:      "CONFIG_PATH",
			Usage:       "optionally, specify the `CONFIG_PATH` environment variable",
		},
		cli.StringFlag{
			Name:        "extra,e",
			Value:       "",
			Destination: &extraArg,
			EnvVar:      "EXTRA_ARG",
			Usage:       "used for testing purposes",
			Hidden:      true,
		},
	}
}

// RunCLI gathers command line arguments and runs the CLI application
func RunCLI(osArgs []string, version, revision string, opts ...runtime.AppOption) error {

	app.Version = fmt.Sprintf("%s (revision %s)", version, revision)

	opts = append(opts, runtime.OsArguments(osArgs))

	rt, err := runtime.New(opts...)
	if err != nil {
		fmt.Printf("unable to initialize CLI: %v", err)
		os.Exit(1)
	}

	//return runtime.Run(args)
	return rt.Run(app)
}
