package provider_init

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
type ProviderInit struct {
	Title            string // Descriptive title - helps humans understand file
	EthereumProvider EthereumProviderInitType
	P2pNetwork       P2pNetworkInitType
	BeaconConfig     BeaconConfigType
}

// This is the information that is read in from the .toml file
type EthereumProviderInitType struct {

	// Example: "ws://192.168.0.157:8546"
	GethConnectionString string

	// Example: "0x6ffba2d0f4c8fd7961f516af43c55fe2d56f6044"
	Account string

	// Names and addresses for contracts that can be called.
	ContractAddress map[string]string

	// Full path to a file like
	// "UTC--2018-02-15T19-57-35.216297214Z--6ffba2d0f4c8fd7961f516af43c55fe2d56f6044"
	KeyFile string

	// Password for accessing the account
	KeyFilePassword string

	// Environment variable to pull KeyFilePassword from environment or if "prompt"
	// then prompt for password.
	EnvKeyFilePassword string
}

// Initialization for P2P network (Just a placeholder until I understand what this is)
type P2pNetworkInitType struct {
	Placeholder string
}

// Initialization for the relay beacon (this may need to be in the on-chain config?)
type BeaconConfigType struct {
	GroupSize int
	Threshold int
}

// ReadConfig reads in the configuration file in .toml format.
func ReadConfig(fn string) (config ProviderInit, err error) {

	if _, err = toml.DecodeFile(fn, &config); err != nil {
		err = fmt.Errorf("Unable to decode .toml file [%s] error [%s]\n", fn, err)
		return
	}

	if config.EthereumProvider.EnvKeyFilePassword == "prompt" {
		config.EthereumProvider.KeyFilePassword = ReadPassword("Enter Account Password: ")
	} else if config.EthereumProvider.EnvKeyFilePassword != "" {
		config.EthereumProvider.KeyFilePassword = os.Getenv(config.EthereumProvider.EnvKeyFilePassword)
	}

	return
}

// ReadPassword prompts a user to enter a password.   The read password uses
// the system password reading call that helps to prevent key loggers from
// capturing the password.
func ReadPassword(prompt string) (password string) {
	fmt.Print(prompt)
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to read password\n")
		return
	}
	return strings.TrimSpace(string(bytePassword))
}
