package maintainer

import "github.com/keep-network/keep-core/pkg/maintainer/wallet"

// Config contains maintainer configuration.
type Config struct {
	// BitcoinDifficulty indicates whether the Bitcoin difficulty maintainer
	// should be started.
	BitcoinDifficulty bool

	// DisableBitcoinDifficultyProxy indicates whether the Bitcoin difficulty
	// maintainer proxy should be disabled. By default, the Bitcoin difficulty
	// maintainer submits the epoch headers to the relay via the proxy. If this
	// option is set to true, the epoch headers will be submitted directly to
	// the relay.
	DisableBitcoinDifficultyProxy bool

	WalletCoordination wallet.Config
	// TODO: Add options for other maintainer tasks, e.g. spv
}
