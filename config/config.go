package config

import (
	"fmt"
	"strings"
	"syscall"

	"github.com/keep-network/keep-core/pkg/chain/ethereum"
	"github.com/keep-network/keep-core/pkg/net/libp2p"
	"github.com/micro/go-config"
	"github.com/micro/go-config/source/env"
	"github.com/micro/go-config/source/file"
	"github.com/micro/go-config/source/flag"
	"golang.org/x/crypto/ssh/terminal"
)

// Environment variables for the Keep client always start with `KEEP`
// Keys are converted to lowercase and split on underscore. Otherwise
// we expect variable names to be the same in the `config.toml` file and
// environment. Here is a list of recognized environment variables:
// KEEP_ETHEREUM_ACCOUNT_KEYFILE_PASSWORD
// KEEP_ETHEREUM_ACCOUNT
// KEEP_ETHEREUM_KEYFILE
// KEEP_ETHEREUM_RANDOM_BEACON_CONTRACT
// KEEP_ETHEREUM_KEEP_GROUP_CONTRACT
// KEEP_ETHEREUM_STAKING_PROXY_CONTRACT
// KEEP_ETHEREUM_URL
// KEEP_ETHEREUM_URL_RPC
// KEEP_LIBP2P_PORT
// KEEP_LIBP2P_PEERS
// KEEP_LIBP2P_SEED

// Config is the top level config structure.
type Config struct {
	Ethereum ethereum.Config
	LibP2P   libp2p.Config
}

type node struct {
	Port                  int
	MyPreferredOutboundIP string
}

type bootstrap struct {
	URLs []string
	Seed int
}

var (
	// KeepOpts contains global application settings
	KeepOpts Config
)

// ReadConfig reads in the configuration file in .toml format.
func ReadConfig(filePath string) (*Config, error) {
	// load configuration from file, environment and options
	configuration := &Config{}
	config.Load(
		// base config from config file
		file.NewSource(file.WithPath(filePath)),
		// override from env variables
		env.NewSource(env.WithStrippedPrefix("KEEP")),
		// override env with flags
		flag.NewSource(),
	)
	// scan the loaded configuration into our struct
	config.Scan(configuration)
	// check values
	fmt.Printf("config.go Password: %s\n", configuration.Ethereum.Account.KeyFilePassword)
	fmt.Printf("config.go Port: %d\n", configuration.LibP2P.Port)
	fmt.Printf("config.go Peers: %s\n", configuration.LibP2P.Peers)
	fmt.Printf("config.go Seed: %d\n", configuration.LibP2P.Seed)

	// get the keyfile password via prompt or environment variable
	if configuration.Ethereum.Account.KeyFilePassword == "prompt" {
		var (
			password string
			err      error
		)
		if password, err = readPassword("Enter Account Password: "); err != nil {
			return nil, err
		}
		configuration.Ethereum.Account.KeyFilePassword = password
	} else {
		if configuration.Ethereum.Account.KeyFilePassword == "" {
			return nil, fmt.Errorf("must supply a valid password to unlock keyfile")
		}
	}

	if configuration.LibP2P.Seed == 0 && len(configuration.LibP2P.Peers) == 0 {
		return nil, fmt.Errorf("either supply a valid bootstrap seed or valid bootstrap URLs")
	}

	if configuration.LibP2P.Seed != 0 && len(configuration.LibP2P.Peers) > 0 {
		return nil, fmt.Errorf("non-bootstrap node should have bootstrap URLs and a seed of 0")
	}

	return configuration, nil
}

// ReadPassword prompts a user to enter a password.   The read password uses
// the system password reading call that helps to prevent key loggers from
// capturing the password.
func readPassword(prompt string) (string, error) {
	fmt.Print(prompt)
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", fmt.Errorf("Unable to read password, error [%s]", err)
	}
	return strings.TrimSpace(string(bytePassword)), nil
}
