package config

import (
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/BurntSushi/toml"
	"github.com/keep-network/keep-core/pkg/chain/ethereum"
	"github.com/keep-network/keep-core/pkg/net/libp2p"
	"github.com/keep-network/keep-core/util"
	"golang.org/x/crypto/ssh/terminal"
)

// Config is the top level config structure.
type Config struct {
	Ethereum ethereum.Config
	Node     libp2p.NodeConfig
}

// DefaultConfig populates the config struc with default values.
func DefaultConfig() *Config {
	cfg := Config{
		Ethereum: ethereum.DefaultConfig,
		Node:     libp2p.DefaultNodeConfig,
	}
	return &cfg
}

// NewConfig returns a struct that has been populated with default values and any values overwritten by the config file.
func NewConfig(filePath string) (Config, error) {
	cfg := DefaultConfig()
	var err error
	*cfg, err = ReadConfig(filePath)
	if err != nil {
		return Config{}, fmt.Errorf("error reading config file: %v", err)
	}
	return *cfg, nil
}

// ValidationError returns validation errors for all config values.
func (c *Config) ValidationError() error {
	var errMsgs []string
	if c.Ethereum.Account.KeyFilePassword == "" {
		errMsgs = append(errMsgs,
			fmt.Sprintf("password is required;  set environment variable (%s) to password or pass 'prompt'",
				passwordEnvVariable))
	}

	var ethCfg = &ethereum.Config{
		URL:               c.Ethereum.URL,
		URLRPC:            c.Ethereum.URLRPC,
		ContractAddresses: c.Ethereum.ContractAddresses,
		Account:           c.Ethereum.Account,
	}
	errMsgs = util.AppendErrMsgs(errMsgs, ethCfg.ValidationError())

	var libp2pCfg = &libp2p.Config{
		NodeConfig: libp2p.NodeConfig{
			Port:  c.Node.Port,
			Seed:  c.Node.Seed,
			Peers: c.Node.Peers,
		},
	}
	errMsgs = util.AppendErrMsgs(errMsgs, libp2pCfg.ValidationError())

	return util.Err(errMsgs)
}

const passwordEnvVariable = "KEEP_ETHEREUM_PASSWORD"

// Opts contains global application settings.
var Opts Config

// ReadConfig reads in the configuration file in .toml format.
//func ReadConfig(filePath string, cfg *Config) error {
// ReadConfig reads in the configuration file in .toml format.
func ReadConfig(filePath string) (Config, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return Config{}, err
	}
	defer util.CloseReadOnlyFile(f)

	cfg := DefaultConfig()
	if _, err := toml.DecodeFile(filePath, &cfg); err != nil {
		return Config{}, err
	}
	envPassword := os.Getenv(passwordEnvVariable)
	if envPassword == "prompt" {
		password, err := readPassword("Enter Account Password: ")
		if err != nil {
			return Config{}, err
		}
		cfg.Ethereum.Account.KeyFilePassword = password
	} else {
		cfg.Ethereum.Account.KeyFilePassword = envPassword
	}
	return *cfg, cfg.ValidationError()
}

// ReadPassword prompts a user to enter a password.   The read password uses the system
// password reading call that helps to prevent key loggers from capturing the password.
func readPassword(prompt string) (string, error) {
	fmt.Print(prompt)
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", fmt.Errorf("unable to read password, error: %v", err)
	}
	return strings.TrimSpace(string(bytePassword)), nil
}
