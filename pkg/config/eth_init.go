package config

import (
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/BurntSushi/toml"
	"golang.org/x/crypto/ssh/terminal"
)

// Top level initialization structure - the data that is read in from the .toml
// file with the --config flag.
// type ProviderInit struct {
type Config struct {
	Ethereum   EthereumConfig
	P2pNetwork LibP2PConfig
}

// environment varialbe with 'prompt' for prompting for password or the password
const EnvPasswordName = "KEEP_ETHERUM_PASSWORD"

type EthereumAccount struct {

	// Example: "0x6ffba2d0f4c8fd7961f516af43c55fe2d56f6044"
	Address string

	// Full path to a file like
	// "UTC--2018-02-15T19-57-35.216297214Z--6ffba2d0f4c8fd7961f516af43c55fe2d56f6044"
	KeyFile string

	// Password for accessing the account
	KeyFilePassword string
}

// This is the information that is read in from the .toml file
type EthereumConfig struct {

	// Example: "ws://192.168.0.157:8546"
	URL string

	// Names and addresses for contracts that can be called or for events received.
	ContractAddresses map[string]string

	Account EthereumAccount
}

// Initialization for P2P network (Just a placeholder until I understand what this is)
type LibP2PConfig struct {
	Placeholder string
}

// ReadConfig reads in the configuration file in .toml format.
func ReadConfig(filePath string) (config Config, err error) {

	if _, err = toml.DecodeFile(filePath, &config); err != nil {
		return config, fmt.Errorf("Unable to decode .toml file [%s] error [%s]\n", filePath, err)
	}

	var pw string
	envPw := os.Getenv(EnvPasswordName)
	if envPw == "prompt" {
		if pw, err = ReadPassword("Enter Account Password: "); err != nil {
			return config, err
		}
		config.Ethereum.Account.KeyFilePassword = pw
	} else if envPw != "" {
		config.Ethereum.Account.KeyFilePassword = envPw
	} else if config.Ethereum.Account.KeyFilePassword != "" {
	} else {
		return config, fmt.Errorf("Password is required.  Set " + EnvPasswordName + " environment variable to password or 'prompt'")
	}

	return config, nil
}

// ReadPassword prompts a user to enter a password.   The read password uses
// the system password reading call that helps to prevent key loggers from
// capturing the password.
func ReadPassword(prompt string) (string, error) {
	fmt.Print(prompt)
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", fmt.Errorf("Unable to read password, error [%s]", err)
	}
	return strings.TrimSpace(string(bytePassword)), nil
}
