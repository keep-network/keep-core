package maintainer

import (
	"context"

	"github.com/keep-network/keep-core/pkg/bitcoin"
)

func Initialize(
	ctx context.Context,
	config Config,
	btcChain bitcoin.Chain,
	chain RelayChain,
) {
	if config.Relay {
		go newRelay(ctx, btcChain, chain)
	}

	// TODO: Launch other maintainer tasks if necessary, e.g. spv. If no task
	//       has been specified - launch all the maintainer tasks.
	// TODO: Cancel all launched tasks if one of the tasks is unable to be
	//       launched, e.g. due to configuration errors.
}
