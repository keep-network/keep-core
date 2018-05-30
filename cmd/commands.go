package cmd

import (
	"fmt"
	"os"
	"path"

	"path/filepath"

	"github.com/urfave/cli"
)

const (
	defaultGroupSize int = 6
	defaultThreshold int = 2

	// DefaultConfigFileName sets default file name; can be changed with --config CLI flag
	DefaultConfigFileName = "config.toml"
)

var (
	// GroupSize indicates the number of members in this relay group
	GroupSize int

	// Threshold indicates the threshold number of members required to perform signature verification
	Threshold int

	// DefaultConfigPath is a config.toml file in the root directory
	DefaultConfigPath = filepath.Join("../", DefaultConfigFileName)

	// Commands contains the list of keep client commands
	Commands = []cli.Command{
		{
			Name:        "validate-config",
			Usage:       "Validates project's configuration file",
			Description: fmt.Sprintf("Validates %s's configuration, which will be in a .toml file", path.Base(os.Args[0])),
			Action:      ValidateConfig,
		},
		{
			Name:        "smoke-test",
			Usage:       "Simulates DKG and signature verification",
			Description: "Simulate Distributed Key Generation (DKG) and verify group's threshold signature",
			Action:      SmokeTest,
			Flags: []cli.Flag{
				&cli.IntFlag{
					Name:        "group-size,g",
					Value:       defaultGroupSize,
					Destination: &GroupSize,
					EnvVar:      "GROUP_SIZE",
					Usage:       "optionally, specify the `GROUP_SIZE` environment variable",
				},
				&cli.IntFlag{
					Name:        "threshold,t",
					Value:       defaultThreshold,
					Destination: &Threshold,
					EnvVar:      "THRESHOLD",
					Usage:       "optionally, specify the `THRESHOLD` environment variable",
				},
			},
		},
	}
)
