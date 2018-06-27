package config

import (
	"fmt"
	"os"
	"strings"
	"syscall"

	"errors"

	"github.com/BurntSushi/toml"
	"github.com/keep-network/keep-core/pkg/chain/ethereum"
	"github.com/keep-network/keep-core/util"
	"golang.org/x/crypto/ssh/terminal"
)

// Err embeds a wrapped error type for using an error number for error type comparisons.
type Err struct {
	util.ErrWrap
}

// Error is a behavior only available for the Err type.
func (e Err) Error() string {
	return fmt.Sprintf("[%d] %s", e.ErrorNumber(), e.Err.Error())
}


// Config is the top level config structure.
type Config struct {
	Ethereum  ethereum.Config
	Bootstrap bootstrap
	Node      node
}

func (c *Config) Error() string {
	return fmt.Sprintf("invalid Ethereum.Account.Address: %s", c.Ethereum.Account.Address)
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

	if cfg.Node.Port == 0 {
		return cfg, fmt.Errorf("missing value for port; see node section in config file or use --port flag")
	}

	if cfg.Bootstrap.Seed == 0 && len(cfg.Bootstrap.URLs) == 0 {
		return cfg, fmt.Errorf("either supply a valid bootstrap seed or valid bootstrap URLs")
	}

	if cfg.Bootstrap.Seed != 0 && len(cfg.Bootstrap.URLs) > 0 {
		return cfg, fmt.Errorf("non-bootstrap node should have bootstrap URLs and a seed of 0")
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
