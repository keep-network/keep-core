package main

// Test that we call the contract and get back a value.

import (
	"flag"
	"fmt"
	"math/big"
	"os"
	"strings"
	"sync"
	"syscall"

	"github.com/BruntSushi/toml"
	"github.com/keep-network/keep-core/pkg/beacon/relay/event"
	"github.com/keep-network/keep-core/pkg/chain/ethereum"
	"github.com/pschlump/MiscLib"
	"github.com/pschlump/godebug"
	"golang.org/x/crypto/ssh/terminal"
)

var fn = flag.String("cfg", "testnet.toml", "Path to configuration file") // 0
var blockReward = flag.Int("blockReward", 1, "Block Reward Value")
var seed = flag.Int("seed", 4, "Random seed value")

func main() {
	flag.Parse()

	fns := flag.Args()
	if len(fns) != 0 {
		fmt.Fprintf(os.Stderr, "Usage: ./RequestRelayEntry [--cfg testnet.toml] [--blockReward 1] [--seed 4]\n")

		os.Exit(1)
	}

	aRun := struct {
		blockReward int64
		seed        int64
	}{
		blockReward: int64(*blockReward),
		seed:        int64(*seed),
	}

	rcfg, err := ReadConfig(*fn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read %s: %s\n", fn, err)
		os.Exit(1)
	}

	// FIXME: This super ugly code is because of a Go dependency loop that I have not had time to resolve.
	cfg := ethereum.Config{
		URL:               rcfg.Ethereum.URL,
		URLRPC:            rcfg.Ethereum.URLRPC,
		ContractAddresses: make(map[string]string),
		Account: ethereum.Account{
			Address:         rcfg.Ethereum.Account.Address,
			KeyFile:         rcfg.Ethereum.Account.KeyFile,
			KeyFilePassword: rcfg.Ethereum.Account.KeyFilePassword,
		},
	}
	for key, val := range rcfg.Ethereum.ContractAddresses {
		cfg.ContractAddresses[key] = val
	}

	if godebug.SVar(rcfg.Ethereum) != godebug.SVar(cfg) {
		fmt.Printf("rcfg: %s\n", godebug.SVarI(rcfg.Ethereum))
		fmt.Printf("cfg: %s\n", godebug.SVarI(cfg))
		fmt.Fprintf(os.Stderr, "Failed to read configuration\n")
		os.Exit(1)
	}

	hdl, err := ethereum.Connect(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to connect: %s\n", err)
		os.Exit(1)
	}

	ec := hdl.ThresholdRelay() // convert to Interface interface

	var wg sync.WaitGroup
	wg.Add(1)

	_ = ec.RequestRelayEntry(
		big.NewInt(aRun.blockReward),
		big.NewInt(aRun.seed),
	).OnSuccess(func(data *event.RelayEntryRequested) {
		fmt.Printf("%ssuccess: %s\n%s", MiscLib.ColorGreen, godebug.SVarI(data), MiscLib.ColorReset)
		wg.Done()
	}).OnFailure(func(err error) {
		fmt.Printf("%serror: %s\n%s", MiscLib.ColorRed, err, MiscLib.ColorReset)
		wg.Done()
	})

	// fmt.Fprintf(os.Stderr, "Waiting for transaction to run (2 minute sleep, sorry...)...\n")
	// time.Sleep(120 * time.Second)
	wg.Wait()
}

// This super ugly code is because of a Go dependency loop that I have not had time to resolve.
// ------------------------------------------------------------------------------------------------------------------------------------
// From config package.
// ------------------------------------------------------------------------------------------------------------------------------------
type rAccount struct {
	Address         string
	KeyFile         string
	KeyFilePassword string
}

type rConfig struct {
	URL               string
	URLRPC            string
	ContractAddresses map[string]string
	Account           rAccount
}

// ------------------------------------------------------------------------------------------------------------------------------------
// From config package.
// ------------------------------------------------------------------------------------------------------------------------------------

const passwordEnvVariable = "KEEP_ETHEREUM_PASSWORD"

// Config is the top level config structure.
type ReadConfigType struct {
	Ethereum  rConfig
	Bootstrap bootstrap
	Node      node
}

type node struct {
	Port                  int
	MyPreferredOutboundIP string
}

type bootstrap struct {
	URLs []string
	Seed int
}

// ReadConfig reads in the configuration file in .toml format.
func ReadConfig(filePath string) (cfg ReadConfigType, err error) {
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
