package maintainer

import (
	"context"

	"github.com/keep-network/keep-core/pkg/bitcoin"
)

func Initialize(
	ctx context.Context,
	config Config,
	btcChain bitcoin.Chain,
	chain BitcoinDifficultyChain,
) {
	// If none of the maintainers was specified in the config (i.e. no option was
	// provided to the `maintainer` command), all maintainers should be launched.
	launchAll := !config.BitcoinDifficulty

	if config.BitcoinDifficulty || launchAll {
		initializeBitcoinDifficultyMaintainer(ctx, btcChain, chain)
	}

	// TODO: Allow for launching multiple maintainers here. Every flag
	//       indicating a maintainer task should launch a separate maintainer.
	//       Notice that panic on one maintainer goroutine will crush the whole
	//       program. Consider cancelling all maintainers if one maintainer
	//       cannot ba launched due to a configuration error.
}
