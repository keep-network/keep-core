package gen

import (
	_ "embed"
	"strings"
)

//go:generate make

var (
	//go:embed _address/RandomBeacon
	randomBeaconAddressFileContent string

	// RandomBeaconAddress is a Random Beacon contract's address read from the NPM package.
	RandomBeaconAddress string = strings.TrimSpace(randomBeaconAddressFileContent)
)
