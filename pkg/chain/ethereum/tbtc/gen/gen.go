package gen

import (
	_ "embed"
	"strings"
)

//go:generate make

var (
	//go:embed _address/Bridge
	bridgeAddressFileContent string

	// BridgeAddress is a Bridge contract's address read from the NPM package.
	BridgeAddress string = strings.TrimSpace(bridgeAddressFileContent)
)
