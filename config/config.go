package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"syscall"

	"github.com/BurntSushi/toml"
	"github.com/keep-network/keep-core/pkg/chain/ethereum"
	"github.com/keep-network/keep-core/pkg/net/libp2p"
	"golang.org/x/crypto/ssh/terminal"
)

const ethPasswordEnvVariable = "KEEP_ETHEREUM_PASSWORD"
const ethAccountEnvVariable = "KEEP_ETHEREUM_ACCOUNT"
const ethKeyfileEnvVariable = "KEEP_ETHEREUM_KEYFILE"
const ethRandomBeaconContractEnvVariable = "KEEP_ETHEREUM_RANDOM_BEACON_CONTRACT"
const ethKeepGroupContractEnvVariable = "KEEP_ETHEREUM_KEEP_GROUP_CONTRACT"
const ethStakingProxyContractEnvVariable = "KEEP_ETHEREUM_STAKING_PROXY_CONTRACT"
const ethURLEnvVariable = "KEEP_ETHEREUM_URL"
const ethURLRPCEnvVariable = "KEEP_ETHEREUM_URL_RPC"
const libp2pPortEnvVariable = "KEEP_LIBP2P_PORT"
const libp2pPeersEnvVariable = "KEEP_LIBP2P_PEERS"
const libp2pSeedEnvVariable = "KEEP_LIBP2P_SEED"

// Config is the top level config structure.
type Config struct {
	Ethereum ethereum.Config
	LibP2P   libp2p.Config
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
func ReadConfig(filePath string) (*Config, error) {
	config := &Config{}
	if _, err := toml.DecodeFile(filePath, config); err != nil {
		return nil, fmt.Errorf("unable to decode .toml file [%s] error [%s]", filePath, err)
	}

	// get the keyfile password via prompt or environment variable
	envPassword := os.Getenv(ethPasswordEnvVariable)
	if envPassword == "prompt" {
		var (
			password string
			err      error
		)
		if password, err = readPassword("Enter Account Password: "); err != nil {
			return nil, err
		}
		config.Ethereum.Account.KeyFilePassword = password
	} else {
		config.Ethereum.Account.KeyFilePassword = envPassword
	}

	if config.Ethereum.Account.KeyFilePassword == "" {
		return nil, fmt.Errorf("Password is required.  Set " + ethPasswordEnvVariable + " environment variable to password or 'prompt'")
	}

	// override account from environment if set or fallback to configfile
	if envAccount, ok := os.LookupEnv(ethAccountEnvVariable); ok {
		config.Ethereum.Account.Address = envAccount
	}
	// complain if account is still not set (not in env or config file)
	if config.Ethereum.Account.Address == "" {
		return nil, fmt.Errorf("Address is required.  Set " + ethAccountEnvVariable + " environment variable to account")
	}

	// get keyfile from environment if set or fallback to config file
	if envKeyfile, ok := os.LookupEnv(ethKeyfileEnvVariable); ok {
		config.Ethereum.Account.KeyFile = envKeyfile
	}
	// complain if keyfile is still not set (not in env or config file)
	if config.Ethereum.Account.KeyFile == "" {
		return nil, fmt.Errorf("Keyfile is required.  Set " + ethKeyfileEnvVariable + " environment variable to keyfile")
	}

	// get random beacon contract address from environment if set or fallback to config file
	if ethRandomBeaconContract, ok := os.LookupEnv(ethRandomBeaconContractEnvVariable); ok {
		config.Ethereum.ContractAddresses["KeepRandomBeacon"] = ethRandomBeaconContract
	}
	// complain if random beacon contract address is still not set (not in env or config file)
	if config.Ethereum.ContractAddresses["KeepRandomBeacon"] == "" {
		return nil, fmt.Errorf("Keep Random Beacon contract address is required.  Set " + ethRandomBeaconContractEnvVariable + " environment variable to Keep Random Beacon contract address.")
	}

	// get Keep group contract address from environment if set or fallback to config file
	if ethKeepGroupContract, ok := os.LookupEnv(ethKeepGroupContractEnvVariable); ok {
		config.Ethereum.ContractAddresses["KeepGroup"] = ethKeepGroupContract
	}
	// complain if Keep group contract address is still not set (not in env or config file)
	if config.Ethereum.ContractAddresses["KeepGroup"] == "" {
		return nil, fmt.Errorf("Keep Group contract address is required.  Set " + ethKeepGroupContractEnvVariable + " environment variable to Keep Group contract address.")
	}

	// get staking proxy contract address from environment if set or fallback to config file
	if ethStakingProxyContract, ok := os.LookupEnv(ethStakingProxyContractEnvVariable); ok {
		config.Ethereum.ContractAddresses["StakingProxy"] = ethStakingProxyContract
	}
	// complain if staking proxy contract address is still not set (not in env or config file)
	if config.Ethereum.ContractAddresses["StakingProxy"] == "" {
		return nil, fmt.Errorf("Keep Random Beacon contract address is required.  Set " + ethStakingProxyContractEnvVariable + " environment variable to Staking Proxy contract address.")
	}

	// get Ethereum URL from environment if set or fallback to config file
	if ethURL, ok := os.LookupEnv(ethURLEnvVariable); ok {
		config.Ethereum.URL = ethURL
	}
	// complain if Ethereum URL is still not set (not in env or config file)
	if config.Ethereum.URL == "" {
		return nil, fmt.Errorf("Keep Ethereum URL is required.  Set " + ethURLEnvVariable + " environment variable to Ethereum URL.")
	}

	// get Ethereum URL RPC from environment if set or fallback to config file
	if ethURLRPC, ok := os.LookupEnv(ethURLRPCEnvVariable); ok {
		config.Ethereum.URLRPC = ethURLRPC
	}
	// complain if Ethereum URL RPC is still not set (not in env or config file)
	if config.Ethereum.URLRPC == "" {
		return nil, fmt.Errorf("Keep Ethereum URL RPC is required.  Set " + ethURLRPCEnvVariable + " environment variable to Ethereum URL RPC.")
	}

	// get Keep LibP2P port from environment if set or fallback to config file
	if libp2pPort, ok := os.LookupEnv(libp2pPortEnvVariable); ok {
		port, err := strconv.Atoi(libp2pPort)
		if err != nil {
			return nil, fmt.Errorf("Error parsing %s: %s", libp2pPortEnvVariable, err)
		}
		config.LibP2P.Port = port
	}

	// Get Keep LibP2P peers from environment if set or fallback to config file.
	// In the toml config file this is defined as a toml serialized string array
	// and we use a similar JSON representation inside the environment variable.
	// The environment variable represents the string array using the following
	// syntax: ["peer0", "peer1", "peer2"]
	if libp2pPeers, ok := os.LookupEnv(libp2pPeersEnvVariable); ok {
		var peers []string
		_ = json.Unmarshal([]byte(libp2pPeers), &peers)
		config.LibP2P.Peers = peers
	}

	// get Keep LibP2P peers from environment if set or fallback to config file
	if libp2pSeed, ok := os.LookupEnv(libp2pSeedEnvVariable); ok {
		seed, err := strconv.Atoi(libp2pSeed)
		if err != nil {
			return nil, fmt.Errorf("Error parsing %s: %s", libp2pSeedEnvVariable, err)
		}
		config.LibP2P.Seed = seed
	}

	if config.LibP2P.Port == 0 {
		return nil, fmt.Errorf("missing value for port; see node section in config file or use --port flag")
	}

	if config.LibP2P.Seed == 0 && len(config.LibP2P.Peers) == 0 {
		return nil, fmt.Errorf("either supply a valid bootstrap seed or valid bootstrap URLs")
	}

	if config.LibP2P.Seed != 0 && len(config.LibP2P.Peers) > 0 {
		return nil, fmt.Errorf("non-bootstrap node should have bootstrap URLs and a seed of 0")
	}

	return config, nil
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
