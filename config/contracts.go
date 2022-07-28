package config

import (
	"fmt"
	"strings"

	"github.com/keep-network/keep-core/pkg/chain/ethereum"

	"github.com/spf13/viper"
)

// GetDeveloperContractAddressKey returns a key for developer contract address configuration.
func GetDeveloperContractAddressKey(contractName string) string {
	return fmt.Sprintf(
		"developer.%sAddress",
		strings.ToLower(contractName[:1])+contractName[1:],
	)
}

// Configuration properties for contract addresses are expected to be defined with
// name `developer.<contractName>Address`. Although the `Config` struct stores them
// in the `Ethereum.ContractAddresses` map. This function binds these two with
// an alias.
func initializeContractAddressesAliases() {
	aliasEthereumContract := func(contractName string) {
		configStructProperty := fmt.Sprintf("Ethereum.ContractAddresses.%s", contractName)
		configKey := GetDeveloperContractAddressKey(contractName)

		viper.RegisterAlias(configStructProperty, configKey)
	}

	aliasEthereumContract(ethereum.RandomBeaconContractName)
	aliasEthereumContract(ethereum.TokenStakingContractName)
	aliasEthereumContract(ethereum.WalletRegistryContractName)
}
