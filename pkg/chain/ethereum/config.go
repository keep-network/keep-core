package ethereum

import (
	"fmt"

	"regexp"

	"github.com/keep-network/keep-core/util"
)

const (
	ethereumURLPattern     = `ws://.+|\w.ipc`
	ethereumURLRPCPattern  = `^https?:\/\/(.+)\.(.+)`
	ethereumAddressPattern = `([13][a-km-zA-HJ-NP-Z1-9]{25,34}|0x[a-fA-F0-9]{40}|\\w+\\.eth(\\W|$)|(?i:iban:)?XE[0-9]{2}[a-zA-Z]{16})|^(bc1|[13])[a-zA-HJ-NP-Z0-9]{25,39}$`
	ethereumKeyfilePattern = `\/.+`
)

// Account is a struct that contains the configuration for accessing an
// Ethereum network and a contract on the network.
type Account struct {
	// Address is the address of this Ethereum account, from which transactions
	// will be sent when interacting with the Ethereum network.
	// Example: "0x6ffba2d0f4c8fd7263f546afaaf25fe2d56f6044".
	Address string

	// Keyfile is a full path to a key file.  Normally this file is one of the
	// imported keys in your local Ethereum server.  It can normally be found in
	// a directory <some-path>/data/keystore/ and starts with its creation date
	// "UTC--.*".
	KeyFile string

	// KeyFilePassword is the password used to unlock the account specified in
	// KeyFile.
	KeyFilePassword string
}

// Config is a struct that contains the configuration needed to connect to an
// Ethereum node.   This information will give access to an Ethereum network.
type Config struct {
	// Example: "ws://192.168.0.157:8546".
	URL string

	// Example: "http://192.168.0.157:8545".
	URLRPC string

	// A  map from contract names to contract addresses.
	ContractAddresses map[string]string

	Account Account
}

// DefaultConfig contains default ethereum package values.
var DefaultConfig = Config{
	URL:    "ws://192.168.0.10:8546",
	URLRPC: "http://192.168.0.10:8545",
	ContractAddresses: map[string]string{
		"KeepRandomBeacon": "0xdeadbeefdeadbeefdeadbeefdeadbeefdeadbeef",
		"GroupContract":    "0xdeadbeefdeadbeefdeadbeefdeadbeefdeadbeef",
	},
	Account: Account{
		Address:         "0xdeadbeefdeadbeefdeadbeefdeadbeefdeadbeef",
		KeyFile:         "0xdeadbeefdeadbeefdeadbeefdeadbeefdeadbeef",
		KeyFilePassword: "not-my-password",
	},
}

// ValidationError returns validation errors for all config values.
func (c *Config) ValidationError() error {
	var errMsgs []string

	r, err := regexp.Compile(ethereumURLPattern)
	if err != nil {
		panic(fmt.Sprintf("Error compiling regexp: %s", ethereumURLPattern))
	}
	if r.FindString(c.URL) == "" {
		errMsgs = append(errMsgs, fmt.Sprintf("Ethereum.URL (%s) invalid; format expected: %s",
			c.URL,
			ethereumURLPattern))
	}

	r, err = regexp.Compile(ethereumURLRPCPattern)
	if err != nil {
		panic(fmt.Sprintf("Error compiling regexp: %s", ethereumURLRPCPattern))
	}
	if r.FindString(c.URLRPC) == "" {
		errMsgs = append(errMsgs, fmt.Sprintf("Ethereum.URLRPC (%s) invalid; format expected: %s",
			c.URLRPC,
			ethereumURLRPCPattern))
	}

	r, err = regexp.Compile(ethereumAddressPattern)
	if err != nil {
		panic(fmt.Sprintf("Error compiling regexp: %s", ethereumAddressPattern))
	}
	if r.FindString(c.Account.Address) == "" {
		errMsgs = append(errMsgs, fmt.Sprintf("Ethereum.Account.Address (%s) invalid; format expected: %s",
			c.Account.Address,
			ethereumAddressPattern))
	}
	if r.FindString(c.Account.KeyFile) == "" {
		errMsgs = append(errMsgs, fmt.Sprintf("Ethereum.Account.KeyFile (%s) invalid; format expected: %s",
			c.Account.KeyFile,
			ethereumKeyfilePattern))
	}
	if r.FindString(c.ContractAddresses["KeepRandomBeacon"]) == "" {
		errMsgs = append(errMsgs, fmt.Sprintf("Ethereum.ContractAddresses[KeepRandomBeacon] (%s) invalid; format expected: %s",
			c.ContractAddresses["KeepRandomBeacon"],
			ethereumURLPattern))
	}
	if r.FindString(c.ContractAddresses["GroupContract"]) == "" {
		errMsgs = append(errMsgs, fmt.Sprintf("Ethereum.ContractAddresses[GroupContract](%s) invalid; format expected: %s",
			c.ContractAddresses["GroupContract"],
			ethereumURLPattern))
	}
	return util.Err(errMsgs)
}
