package config

import (
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

const (
	passwordEnvVariable    = "KEEP_ETHEREUM_PASSWORD"
	ethereumURLPattern     = `ws://.+|\w.ipc`
	ethereumURLRPCPattern  = `^https?:\/\/(.+)\.(.+)`
	ethereumAddressPattern = `([13][a-km-zA-HJ-NP-Z1-9]{25,34}|0x[a-fA-F0-9]{40}|\\w+\\.eth(\\W|$)|(?i:iban:)?XE[0-9]{2}[a-zA-Z]{16})|^(bc1|[13])[a-zA-HJ-NP-Z0-9]{25,39}$`
	ethereumKeyfilePattern = `\/.+`
)

var (
	// Opts contains global application settings.
	Opts            Config
	ethURLRegex     = util.CompileRegex(ethereumURLPattern)
	ethURLRPCRegex  = util.CompileRegex(ethereumURLRPCPattern)
	ethAddressRegex = util.CompileRegex(ethereumAddressPattern)
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
