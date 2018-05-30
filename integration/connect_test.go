package integration

import (
	"flag"
	"fmt"
	"os"
	"testing"

	"github.com/keep-network/keep-core/pkg/chain/ethereum"
	"github.com/keep-network/keep-core/pkg/config"
)

var client *ethereum.Provider

func TestConnectToServer(t *testing.T) {
}

var (
	integration = flag.Bool("integration", false, "run Geth integration tests")
	configFile  = flag.String("config", "test/config.toml", "Config file to specify connection to ethereum")
)

func TestMain(m *testing.M) {

	flag.Parse()

	os.Setenv("KEEP_ETHEREUM_PASSWORD", "not-my-password")

	Config, err := config.ReadConfig(*configFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "FAIL: Error reading configuration: %s\n", err)
		os.Exit(1)
	}

	if *integration {
		client, err = ethereum.Connect(Config.Ethereum)
		if err != nil {
			fmt.Fprintf(os.Stderr, "FAIL: Failed to connect to Ethereum: %s\n", err)
			os.Exit(1)
		}

		os.Exit(m.Run())
	}
}
