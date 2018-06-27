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

	// EthereumURLErrNo is the error code for an invalid EthereumURL.
	EthereumURLErrNo
	// EthereumURLRPCErrNo is the error code for an invalid EthereumURLRPC.
	EthereumURLRPCErrNo
	// EthereumAccountAddressErrNo is the error code for an invalid EthereumAccountAddress.
	EthereumAccountAddressErrNo
	// EthereumAccountKeyfileErrNo is the error code for an invalid EthereumAccountKeyfile.
	EthereumAccountKeyfileErrNo
	// EthereumContractKeepRandomBeaconAddressErrNo is the error code for an invalid EthereumContractKeepRandomBeaconAddress.
	EthereumContractKeepRandomBeaconAddressErrNo
	// EthereumContractGroupContractAddressErrNo is the error code for an invalid EthereumContractGroupContractAddress.
	EthereumContractGroupContractAddressErrNo
	// NodePortErrNo is the error code for an invalid NodePort.
	NodePortErrNo
	// BootstrapSeedAndURLsErrNo is the error code for an invalid BootstrapSeedAndURLs.
	BootstrapSeedAndURLsErrNo
	// PeerURLsAndSeedErrNo is the error code for an invalid PeerURLsAndSeed.
	PeerURLsAndSeedErrNo

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

	ethereumURLPattern     = `ws://.+|\w.ipc`
	ethereumURLRPCPattern  = `^https?:\/\/(.+)\.(.+)`
	ethereumAddressPattern = `([13][a-km-zA-HJ-NP-Z1-9]{25,34}|0x[a-fA-F0-9]{40}|\\w+\\.eth(\\W|$)|(?i:iban:)?XE[0-9]{2}[a-zA-Z]{16})|^(bc1|[13])[a-zA-HJ-NP-Z0-9]{25,39}$`
	ethereumKeyfilePattern = `\/.+`
)

var (
	// KeepOpts contains global application settings
	KeepOpts Config
	ethURLRegex     = util.CompileRegex(ethereumURLPattern)
	ethURLRPCRegex  = util.CompileRegex(ethereumURLRPCPattern)
	ethAddressRegex = util.CompileRegex(ethereumAddressPattern)
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
	if !util.MatchFound(ethURLRegex, cfg.Ethereum.URL) {
		err := fmt.Errorf("Ethereum.URL (%s) invalid; format expected:%s",
			cfg.Ethereum.URL, ethereumURLPattern)

		return cfg, util.ErrWrap{ErrNo: EthereumURLErrNo, Err: err}
	}

	if !util.MatchFound(ethURLRPCRegex, cfg.Ethereum.URLRPC) {
		err := fmt.Errorf("Ethereum.URLRPC (%s) invalid; format expected:%s",
			cfg.Ethereum.URLRPC, ethereumURLRPCPattern)

		return cfg, util.ErrWrap{ErrNo: EthereumURLRPCErrNo, Err: err}
	}

	if !util.MatchFound(ethAddressRegex, cfg.Ethereum.Account.Address) {
		err := fmt.Errorf("Ethereum.Account.Address (%s) invalid; format expected:%s",
			cfg.Ethereum.Account.Address, ethereumAddressPattern)

		return cfg, util.ErrWrap{ErrNo: EthereumAccountAddressErrNo, Err: err}
	}

	if !util.MatchFound(ethAddressRegex, cfg.Ethereum.Account.KeyFile) {
		err := fmt.Errorf("Ethereum.Account.KeyFile (%s) invalid; format expected:%s",
			cfg.Ethereum.Account.KeyFile, ethereumKeyfilePattern)

		return cfg, util.ErrWrap{ErrNo: EthereumAccountKeyfileErrNo, Err: err}
	}

	if !util.MatchFound(ethAddressRegex, cfg.Ethereum.ContractAddresses["KeepRandomBeacon"]) {
		err := fmt.Errorf("Ethereum.ContractAddresses[KeepRandomBeacon] (%s) invalid; format expected:%s",
			cfg.Ethereum.ContractAddresses["KeepRandomBeacon"], ethereumAddressPattern)

		return cfg, util.ErrWrap{ErrNo: EthereumContractKeepRandomBeaconAddressErrNo, Err: err}
	}

	if !util.MatchFound(ethAddressRegex, cfg.Ethereum.ContractAddresses["GroupContract"]) {
		err := fmt.Errorf("Ethereum.ContractAddresses[GroupContract] (%s) invalid; format expected:%s",
			cfg.Ethereum.ContractAddresses["GroupContract"], ethereumAddressPattern)

		return cfg, util.ErrWrap{ErrNo: EthereumContractGroupContractAddressErrNo, Err: err}
	}

	if cfg.Node.Port == 0 {
		err := errors.New("missing value for port; see node section in config file or use --port flag")
		return cfg, util.ErrWrap{ErrNo: NodePortErrNo, Err: err}
	}

	if len(cfg.Bootstrap.URLs) == 0 && cfg.Bootstrap.Seed <= 0 {
		err := errors.New("either supply a valid bootstrap seed or valid bootstrap URLs")
		return cfg, util.ErrWrap{ErrNo: BootstrapSeedAndURLsErrNo, Err: err}
	}

	if len(cfg.Bootstrap.URLs) > 0 && cfg.Bootstrap.Seed != 0 {
		err := errors.New("non-bootstrap node should have bootstrap URLs and a seed of 0")
		return cfg, util.ErrWrap{ErrNo: PeerURLsAndSeedErrNo, Err: err}
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
