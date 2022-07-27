package gen

import (
	_ "embed"
)

//go:generate make download_artifacts
//go:generate make

//go:embed _address/TokenStaking
var TokenStakingAddress string
