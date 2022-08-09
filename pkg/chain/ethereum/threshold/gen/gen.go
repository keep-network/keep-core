package gen

import (
	_ "embed"
	"strings"
)

//go:generate make

var (
	//go:embed _address/TokenStaking
	tokenStakingAddressFileContent string

	// TokenStakingAddress is a TokenStaking contract's address read from the NPM package.
	TokenStakingAddress string = strings.TrimSpace(tokenStakingAddressFileContent)
)
