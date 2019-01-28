package config

import (
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/BurntSushi/toml"
	"github.com/keep-network/keep-core/pkg/chain/ethereum"
	"github.com/keep-network/keep-core/pkg/net/libp2p"
	"golang.org/x/crypto/ssh/terminal"
)

const ethPasswordEnvVariable = "KEEP_ETHEREUM_PASSWORD"
const ethAccountEnvVariable = "KEEP_ETHEREUM_ACCOUNT"
const ethKeyfileEnvVariable = "KEEP_ETHEREUM_KEYFILE"

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
	config := &Config{}
	if _, err := toml.DecodeFile(filePath, config); err != nil {
		return nil, fmt.Errorf("unable to decode .toml file [%s] error [%s]", filePath, err)
	}

	// get the keyfile password via prompt or environment variable
	envPassword := os.Getenv(ethPasswordEnvVariable)
	if envPassword == "prompt" {
		var (
			password string
			err      error
		)
		if password, err = readPassword("Enter Account Password: "); err != nil {
			return nil, err
		}
		config.Ethereum.Account.KeyFilePassword = password
	} else {
		config.Ethereum.Account.KeyFilePassword = envPassword
	}

	if config.Ethereum.Account.KeyFilePassword == "" {
		return nil, fmt.Errorf("Password is required.  Set " + ethPasswordEnvVariable + " environment variable to password or 'prompt'")
	}

	// override account from environment if set or fallback to configfile
	if envAccount, ok := os.LookupEnv(ethAccountEnvVariable); ok {
		config.Ethereum.Account.Address = envAccount
	}
	// complain if account is still not set (not in env or config file)
	if config.Ethereum.Account.Address == "" {
		return nil, fmt.Errorf("Address is required.  Set " + ethAccountEnvVariable + " environment variable to account")
	}

	// get keyfile from environment if set or fallback to config file
	if envKeyfile, ok := os.LookupEnv(ethKeyfileEnvVariable); ok {
		config.Ethereum.Account.KeyFile = envKeyfile
	}
	// complain if keyfile is still not set (not in env or config file)
	if config.Ethereum.Account.KeyFile == "" {
		return nil, fmt.Errorf("Keyfile is required.  Set " + ethKeyfileEnvVariable + " environment variable to keyfile")
	}

	if config.LibP2P.Port == 0 {
		return nil, fmt.Errorf("missing value for port; see node section in config file or use --port flag")
	}

	if config.LibP2P.Seed == 0 && len(config.LibP2P.Peers) == 0 {
		return nil, fmt.Errorf("either supply a valid bootstrap seed or valid bootstrap URLs")
	}

	if config.LibP2P.Seed != 0 && len(config.LibP2P.Peers) > 0 {
		return nil, fmt.Errorf("non-bootstrap node should have bootstrap URLs and a seed of 0")
	}

	return config, nil
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
