package main

import (
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/BurntSushi/toml"
	"github.com/keep-network/keep-core/pkg/chain/ethereum"
	"golang.org/x/crypto/ssh/terminal"
)

// environment varialbe with 'prompt' for prompting for password or the password
const envPasswordName = "KEEP_ETHERUM_PASSWORD"

// Top level initialization structure - the data that is read in from the .toml
// file with the `--config <file>` flag.
type config struct {
	Ethereum ethereum.EthereumConfig
	// TODO: add this when it is known: P2pNetwork LibP2PConfig
}

// ReadConfig reads in the configuration file in .toml format.
func readConfig(filePath string) (cfg config, err error) {

	if _, err = toml.DecodeFile(filePath, &cfg); err != nil {
		return cfg, fmt.Errorf("Unable to decode .toml file [%s] error [%s]\n", filePath, err)
	}

	var password string
	envPassword := os.Getenv(envPasswordName)
	if envPassword == "prompt" {
		if password, err = readPassword("Enter Account Password: "); err != nil {
			return cfg, err
		}
		cfg.Ethereum.Account.KeyFilePassword = password
	} else {
		cfg.Ethereum.Account.KeyFilePassword = envPassword
	}

	if cfg.Ethereum.Account.KeyFilePassword == "" {
		return cfg, fmt.Errorf("Password is required.  Set " + envPasswordName + " environment variable to password or 'prompt'")
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
