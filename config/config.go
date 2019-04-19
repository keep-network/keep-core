package config

import (
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/keep-network/keep-core/pkg/chain/ethereum"
	"github.com/keep-network/keep-core/pkg/net/libp2p"
	"github.com/micro/go-config"
	"github.com/micro/go-config/source/consul"
	"github.com/micro/go-config/source/file"
	"golang.org/x/crypto/ssh/terminal"
)

const ethPasswordEnvVariable = "KEEP_ETHEREUM_ACCOUNT_KEYFILEPASSWORD"

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
// If consulServer is specified it will try to read configurations
// from a Consul server instead.
func ReadConfig(filePath string, consulServer string) (*Config, error) {
	configuration := &Config{}
	if consulServer != "" {
		fmt.Printf("consulServer set: %s\n", consulServer)
		if err := config.Load(consul.NewSource(consul.WithAddress(consulServer),
			consul.WithPrefix(""))); err != nil {
			return nil, fmt.Errorf(
				"unable to connect to Consul server [%s] error [%s]", consulServer, err)
		}
	} else {
		fmt.Printf("filePath set: %s\n", filePath)
		if err := config.Load(file.NewSource(file.WithPath(filePath))); err != nil {
			return nil, fmt.Errorf(
				"unable to decode .toml file [%s] error [%s]", filePath, err)
		}
	}

	// scan the loaded configuration into our struct
	config.Scan(configuration)

	//DEBUG: create a map of the configuration
	// cMap := config.Map()
	// fmt.Println("map: ", cMap)
	// fmt.Println("Ethereum.Account: ", configuration.Ethereum.Account)
	// fmt.Println("LibP2P.Port: ", configuration.LibP2P.Port)
	// fmt.Println("LibP2P.Seed: ", configuration.LibP2P.Seed)
	// fmt.Println("LibP2P.Peers: ", configuration.LibP2P.Peers)
	fmt.Println("Ethereum.ContractAddresses: ", configuration.Ethereum.ContractAddresses)

	envPassword := os.Getenv(ethPasswordEnvVariable)
	if envPassword == "prompt" {
		var (
			password string
			err      error
		)
		if password, err = readPassword("Enter Account Password: "); err != nil {
			return nil, err
		}
		configuration.Ethereum.Account.KeyFilePassword = password
	} else {
		configuration.Ethereum.Account.KeyFilePassword = envPassword
	}

	if configuration.Ethereum.Account.KeyFilePassword == "" {
		return nil, fmt.Errorf("Password is required.  Set " + ethPasswordEnvVariable + " environment variable to password or 'prompt'")
	}

	if configuration.LibP2P.Port == 0 {
		return nil, fmt.Errorf("missing value for port; see node section in config file or use --port flag")
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
