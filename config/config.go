package config

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/BurntSushi/toml"
	"github.com/keep-network/keep-core/pkg/chain/ethereum"
	"github.com/keep-network/keep-core/util"
	"golang.org/x/crypto/ssh/terminal"
)

// Config is the top level config structure.
type Config struct {
	Ethereum  ethereum.Config
	Bootstrap bootstrap
	Node      node
}

// ValidationError returns validation errors for all config values
func (c *Config) ValidationError() error {
	var errMsgs []string
	if c.Ethereum.Account.KeyFilePassword == "" {
		errMsgs = append(errMsgs,
			fmt.Sprintf("password is required;  set environment variable (%s) to password or pass 'prompt'",
				passwordEnvVariable))
	}
	if !util.MatchFound(ethURLRegex, c.Ethereum.URL) {
		errMsgs = append(errMsgs, fmt.Sprintf("Ethereum.URL (%s) invalid; format expected: %s",
			c.Ethereum.URL,
			ethereumURLPattern))
	}
	if !util.MatchFound(ethURLRPCRegex, c.Ethereum.URLRPC) {
		errMsgs = append(errMsgs, fmt.Sprintf("Ethereum.URLRPC (%s) invalid; format expected: %s",
			c.Ethereum.URLRPC,
			ethereumURLRPCPattern))
	}
	if !util.MatchFound(ethAddressRegex, c.Ethereum.Account.Address) {
		errMsgs = append(errMsgs, fmt.Sprintf("Ethereum.Account.Address (%s) invalid; format expected: %s",
			c.Ethereum.Account.Address,
			ethereumAddressPattern))
	}
	if !util.MatchFound(ethAddressRegex, c.Ethereum.Account.KeyFile) {
		errMsgs = append(errMsgs, fmt.Sprintf("Ethereum.Account.KeyFile (%s) invalid; format expected: %s",
			c.Ethereum.Account.KeyFile,
			ethereumKeyfilePattern))
	}
	if !util.MatchFound(ethAddressRegex, c.Ethereum.ContractAddresses["KeepRandomBeacon"]) {
		errMsgs = append(errMsgs,
			fmt.Sprintf("Ethereum.ContractAddresses[KeepRandomBeacon] (%s) invalid; format expected: %s",
				c.Ethereum.ContractAddresses["KeepRandomBeacon"],
				ethereumAddressPattern))
	}
	if !util.MatchFound(ethAddressRegex, c.Ethereum.ContractAddresses["GroupContract"]) {
		errMsgs = append(errMsgs,
			fmt.Sprintf("Ethereum.ContractAddresses[GroupContract] (%s) invalid; format expected: %s",
				c.Ethereum.ContractAddresses["GroupContract"],
				ethereumAddressPattern))
	}
	if c.Node.Port <= 0 {
		errMsgs = append(errMsgs,
			fmt.Sprintf("Node.Port (%d) invalid; see node section in config file or use --port flag",
				c.Node.Port))
	}
	if len(c.Bootstrap.URLs) == 0 && c.Bootstrap.Seed <= 0 {
		errMsgs = append(errMsgs, fmt.Sprintf("either supply valid Bootstrap.URLs or a Bootstrap.Seed"))
	}
	if len(c.Bootstrap.URLs) > 0 && c.Bootstrap.Seed != 0 {
		errMsgs = append(errMsgs, fmt.Sprintf("non-bootstrap node should have Bootstrap.URL and a Bootstrap.Seed of 0"))
	}
	if len(c.Bootstrap.URLs) > 0 {
		for _, ipfsURL := range c.Bootstrap.URLs {
			if !util.MatchFound(ifpsURLRegex, ipfsURL) {
				errMsgs = append(errMsgs,
					fmt.Sprintf("Bootstrap.URL (%s) invalid; format expected: %s",
						ipfsURL,
						ipfsURLPattern))
			}
		}
	}
	var err error
	if len(errMsgs) > 0 {
		err = errors.New(strings.Join(errMsgs[:], "\n"))
	}
	return err
}

type node struct {
	Port                  int
	MyPreferredOutboundIP string
}

type bootstrap struct {
	URLs []string
	Seed int
}

const (
	passwordEnvVariable    = "KEEP_ETHEREUM_PASSWORD"
	ethereumURLPattern     = `ws://.+|\w.ipc`
	ethereumURLRPCPattern  = `^https?:\/\/(.+)\.(.+)`
	ethereumAddressPattern = `([13][a-km-zA-HJ-NP-Z1-9]{25,34}|0x[a-fA-F0-9]{40}|\\w+\\.eth(\\W|$)|(?i:iban:)?XE[0-9]{2}[a-zA-Z]{16})|^(bc1|[13])[a-zA-HJ-NP-Z0-9]{25,39}$`
	ethereumKeyfilePattern = `\/.+`
	ipfsURLPattern         = `.+\/.*`
)

var (
	// Opts contains global application settings.
	Opts            Config
	ethURLRegex     = util.CompileRegex(ethereumURLPattern)
	ethURLRPCRegex  = util.CompileRegex(ethereumURLRPCPattern)
	ethAddressRegex = util.CompileRegex(ethereumAddressPattern)
	ifpsURLRegex    = util.CompileRegex(ipfsURLPattern)
)

// ReadConfig reads in the configuration file in .toml format.
func ReadConfig(filePath string) (Config, error) {
	var cfg Config
	if _, err := toml.DecodeFile(filePath, &cfg); err != nil {
		return cfg, err
	}
	envPassword := os.Getenv(passwordEnvVariable)
	if envPassword == "prompt" {
		password, err := readPassword("Enter Account Password: ")
		if err != nil {
			return cfg, err
		}
		cfg.Ethereum.Account.KeyFilePassword = password
	} else {
		cfg.Ethereum.Account.KeyFilePassword = envPassword
	}
	return cfg, cfg.ValidationError()
}

// ReadPassword prompts a user to enter a password.   The read password uses the system
// password reading call that helps to prevent key loggers from capturing the password.
func readPassword(prompt string) (string, error) {
	fmt.Print(prompt)
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", fmt.Errorf("Unable to read password, error [%s]", err)
	}
	return strings.TrimSpace(string(bytePassword)), nil
}
