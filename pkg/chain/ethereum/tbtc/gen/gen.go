package gen

import (
	_ "embed"
)

//go:generate make download_artifacts
//go:generate make

//go:embed _address/Bridge
var BridgeAddress string
