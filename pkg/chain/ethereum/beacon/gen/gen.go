package gen

import (
	_ "embed"
)

//go:generate make download_artifacts
//go:generate make

//go:embed _address/RandomBeacon
var RandomBeaconAddress string
