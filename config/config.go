package config

import (
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/BurntSushi/toml"
	"github.com/keep-network/keep-core/pkg/chain/ethereum"
	"golang.org/x/crypto/ssh/terminal"
)

const passwordEnvVariable = "KEEP_ETHEREUM_PASSWORD"

// Config is the top level config structure.
type Config struct {
	Ethereum  ethereum.Config
	Bootstrap bootstrap
	Node  node
}

type node struct {
	Port        int
}

type bootstrap struct {
	URLs        []string
	Seed        int
}

var (
	// KeepOpts contains global application settings
	KeepOpts Config
)

// ReadConfig reads in the configuration file in .toml format.
func ReadConfig(filePath string) (cfg Config, err error) {

	if _, err = toml.DecodeFile(filePath, &cfg); err != nil {
		return cfg, fmt.Errorf("unable to decode .toml file [%s] error [%s]", filePath, err)
	}

	var password string
	envPassword := os.Getenv(passwordEnvVariable)
	if envPassword == "prompt" {
		if password, err = readPassword("Enter Account Password: "); err != nil {
			return cfg, err
		}
		cfg.Ethereum.Account.KeyFilePassword = password
	} else {
		cfg.Ethereum.Account.KeyFilePassword = envPassword
	}

	if cfg.Ethereum.Account.KeyFilePassword == "" {
		return cfg, fmt.Errorf("Password is required.  Set " + passwordEnvVariable + " environment variable to password or 'prompt'")
	}

	if len(cfg.Bootstrap.URLs) == 0 {
		return cfg, fmt.Errorf("Bootstrap URLs missing")
	}

	if cfg.Node.Port == 0 {
		return cfg, fmt.Errorf("Node Port seed missing")
	}

	return cfg, nil
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
