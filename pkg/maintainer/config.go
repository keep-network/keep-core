package maintainer

import (
	"github.com/keep-network/keep-core/pkg/maintainer/btcdiff"
	"github.com/keep-network/keep-core/pkg/maintainer/spv"
)

// Config contains maintainer configuration.
type Config struct {
	BitcoinDifficulty btcdiff.Config
	Spv               spv.Config
}
