package maintainer

import (
	"context"

	"github.com/ipfs/go-log/v2"

	"github.com/keep-network/keep-core/pkg/bitcoin"
	"github.com/keep-network/keep-core/pkg/maintainer/btcdiff"
	"github.com/keep-network/keep-core/pkg/maintainer/spv"
	"github.com/keep-network/keep-core/pkg/maintainer/wallet"
)

var logger = log.Logger("keep-maintainer")

func Initialize(
	ctx context.Context,
	config Config,
	btcChain bitcoin.Chain,
	btcDiffChain btcdiff.Chain,
	coordinatorChain wallet.Chain,
	spvChain spv.Chain,
) {
	// If none of the maintainers was specified in the config (i.e. no option was
	// provided to the `maintainer` command), all maintainers should be launched.
	launchAll := !config.BitcoinDifficulty.Enabled &&
		!config.WalletCoordination.Enabled &&
		!config.Spv.Enabled

	if launchAll {
		logger.Info("initializing all maintainer modules...")
	}

	if config.BitcoinDifficulty.Enabled || launchAll {
		btcdiff.Initialize(
			ctx,
			config.BitcoinDifficulty,
			btcChain,
			btcDiffChain,
		)
	}

	if config.WalletCoordination.Enabled || launchAll {
		wallet.Initialize(
			ctx,
			config.WalletCoordination,
			coordinatorChain,
			btcChain,
		)
	}

	if config.Spv.Enabled || launchAll {
		spv.Initialize(
			ctx,
			config.Spv,
			spvChain,
			btcDiffChain,
			btcChain,
		)
	}

	// TODO: Allow for launching multiple maintainers here. Every flag
	//       indicating a maintainer task should launch a separate maintainer.
	//       Notice that panic on one maintainer goroutine will crush the whole
	//       program. Consider cancelling all maintainers if one maintainer
	//       cannot ba launched due to a configuration error.
}
