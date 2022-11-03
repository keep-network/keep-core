package maintainer

import (
	"context"
)

func Initialize(config Config, ctx context.Context) error {
	// TODO: Add connection to the Bitcoin chain:
	// btcChain, err := bitcoin.Connect(ctx, &maintainerConfig.Bitcoin)
	// if err != nil {
	// 	return fmt.Errorf("could not connect BTC chain: [%v]", err)
	// }

	// TODO: Add connection to the tbtc chain:
	// relayChain, err := connectRelayChain(config)
	// if err != nil {
	// 	return fmt.Errorf("could not connect relay chain: [%v]", err)
	// }

	if config.Relay {
		go newRelay(ctx, nil, nil)
	}

	// TODO: Launch other maintainer tasks if necessary, e.g. spv. If no task
	//       has been specified - launch all the maintainer tasks.
	// TODO: Cancel all launched tasks if one of the tasks is unable to be
	//       launched, e.g. due to configuration errors.

	return nil
}
