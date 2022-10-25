package gen

import (
	_ "embed"
	"strings"
)

// FIXME: Commented out temporarily for mainnet build.
//go:generate make

var (
	// FIXME: Commented out temporarily for mainnet build.
	//go:embed _address/Bridge
	bridgeAddressFileContent string

	// BridgeAddress is a Bridge contract's address read from the NPM package.
	BridgeAddress string = strings.TrimSpace(bridgeAddressFileContent)
)
