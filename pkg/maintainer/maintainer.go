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

	// TODO: Check if launching multiple maintainers is really beneficial.
	//       Panic on one maintainer goroutine will crush the whole program.
	//       Also, if we misconfigure one maintainer and it cannot launch,
	//       should we cancel all the maintainers that already launched properly?

	// TODO: If we decide to enable launching multiple maintainers with one
	//       command, allow them to be launch all together with
	//       `./keep-client maintainer` command (without providing any flag).
}
