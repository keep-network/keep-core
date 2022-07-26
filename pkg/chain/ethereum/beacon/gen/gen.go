package gen

import (
	_ "embed"
	"strings"
)

//go:generate make download_artifacts
//go:generate make

var (
	//go:embed _address/RandomBeacon
	randomBeaconAddressFileContent string

	// RandomBeaconAddress is a WalletRegistry contract's address read from the NPM package.
	RandomBeaconAddress string = strings.TrimSpace(randomBeaconAddressFileContent)
)
