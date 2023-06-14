package maintainer

import (
	"context"

	"github.com/ipfs/go-log/v2"

	"github.com/keep-network/keep-core/pkg/bitcoin"
	"github.com/keep-network/keep-core/pkg/maintainer/wallet"
)

var logger = log.Logger("keep-maintainer")

func Initialize(
	ctx context.Context,
	config Config,
	btcChain bitcoin.Chain,
	btcDiffChain BitcoinDifficultyChain,
	coordinatorChain wallet.Chain,
) {
	// If none of the maintainers was specified in the config (i.e. no option was
	// provided to the `maintainer` command), all maintainers should be launched.
	launchAll := !config.BitcoinDifficulty && !config.WalletCoordination

	if config.BitcoinDifficulty || launchAll {
		initializeBitcoinDifficultyMaintainer(
			ctx,
			btcChain,
			btcDiffChain,
			config.DisableBitcoinDifficultyProxy,
			bitcoinDifficultyDefaultIdleBackOffTime,
			bitcoinDifficultyDefaultRestartBackoffTime,
		)
	}

	if config.WalletCoordination || launchAll {
		wallet.Initialize(
			ctx,
			coordinatorChain,
			btcChain,
		)
	}

	// TODO: Allow for launching multiple maintainers here. Every flag
	//       indicating a maintainer task should launch a separate maintainer.
	//       Notice that panic on one maintainer goroutine will crush the whole
	//       program. Consider cancelling all maintainers if one maintainer
	//       cannot ba launched due to a configuration error.
}
