package maintainer

import (
	"github.com/keep-network/keep-core/pkg/maintainer/btcdiff"
	"github.com/keep-network/keep-core/pkg/maintainer/wallet"
)

// Config contains maintainer configuration.
type Config struct {
	BitcoinDifficulty  btcdiff.Config
	WalletCoordination wallet.Config
}
