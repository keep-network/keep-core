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
	if config.BitcoinDifficulty {
		initializeBitcoinDifficultyMaintainer(ctx, btcChain, chain)
	}

	// TODO: Allow for launching multiple maintainers here. Every flag
	//       indicating a maintainer task should launch a separate maintainer.
	//       If there are no flags specified, all the maintainers should be launch.
	//       Notice that panic on one maintainer goroutine will crush the whole
	//       program. Consider cancelling all maintainers if one maintainer
	//       cannot ba launched due to a configuration error.
}
