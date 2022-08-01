package gen

import (
	_ "embed"
	"strings"
)

//go:generate make download_artifacts
//go:generate make

var (
	//go:embed _address/WalletRegistry
	walletRegistryAddressFileContent string

	// WalletRegistryAddress is a WalletRegistry contract's address read from the NPM package.
	WalletRegistryAddress string = strings.TrimSpace(walletRegistryAddressFileContent)
)
