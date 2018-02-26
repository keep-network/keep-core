package cli

import (
	"flag"
	"fmt"
	"os"
	"strings"

	EventLoop "github.com/keep-network/keep-core/go/interface/MainEventLoop"
)

type CLI struct {
	Args []string
}

func (cli *CLI) Run(args ...string) {

	// Optional set of arguments can be passed
	if len(args) > 0 {
		cli.Args = args // Used for automatic testing purposes
	} else {
		cli.Args = os.Args // Pull in command line arguments
	}

	// cfgLib.ReadConfigFile("cfg.json", &g_cfg)

	// Create a go-like command structure
	switch cli.Args[1] {
	case "version":
		cli.version()

	case "status":
		cmd := flag.NewFlagSet("status", flag.ExitOnError)
		debugFlags := cmd.String("debug", "", "Debug flags")
		cmd.StringVar(debugFlags, "D", "", "Debug flags")
		err := cmd.Parse(cli.Args[2:])
		FatalIfError(err)
		if cmd.Parsed() {
			if *debugFlags != "" {
				SetDebugFlags(*debugFlags)
			}
			cli.status()
		}
	case "startRelay", "relayStart":
		cmd := flag.NewFlagSet("status", flag.ExitOnError)
		debugFlags := cmd.String("debug", "", "Debug flags")
		cmd.StringVar(debugFlags, "D", "", "Debug flags")
		cfgFile := cmd.String("cfg", "", "Config file")
		cmd.StringVar(cfgFile, "c", "", "Config file")
		err := cmd.Parse(cli.Args[2:])
		FatalIfError(err)
		if cmd.Parsed() {
			if *debugFlags != "" {
				SetDebugFlags(*debugFlags)
			}
			cli.relayStart(*cfgFile)
		}

	default:
		usage()
	}

	os.Exit(0)

}

func usage() {
	fmt.Printf(`usage: %s <command> <options>
`)
	os.Exit(1)
}

func SetDebugFlags(flags string) {
	dbs := strings.Split(flags, ",")
	for _, db := range dbs {
		if !Db[db] { // toggle, if not defined then define it.
			Db[db] = true
		} else {
			Db[db] = false
		}
	}
}

var Db map[string]bool

func init() {
	Db = make(map[string]bool)
	// Db["fx1"] = false // Output from processing of functions like __include__
	Db["show-mining"] = true
}

func FatalIfError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal: error %s\n", err)
		os.Exit(1)
	}
}

func (cli *CLI) version() {
	fmt.Printf("Version: KeepRelay 0.0.1\n")
}

func (cli *CLI) status() {
	fmt.Printf("This should printout the satus of the Keep system\n")
}

func (cli *CLI) relayStart(cfgFileName string) {
	// TODO: open config file
	EventLoop.MainEventLoop()
}
