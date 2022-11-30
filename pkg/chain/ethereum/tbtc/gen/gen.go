package gen

import (
	_ "embed"
	"strings"
)

//go:generate make

var (
	bridgeAddressFileContent string

	// BridgeAddress is a Bridge contract's address read from the NPM package.
	BridgeAddress string = strings.TrimSpace(bridgeAddressFileContent)

	//go:embed _address/LightRelay
	lightRelayAddressFileContent string

	// LightRelayAddress is a LightRelay contract's address read from the NPM
	// package.
	LightRelayAddress string = strings.TrimSpace(lightRelayAddressFileContent)
)
