package maintainer

import (
	"context"
	"fmt"
	"time"

	"github.com/ipfs/go-log"
	"github.com/keep-network/keep-core/pkg/bitcoin"
)

var logger = log.Logger("keep-maintainer-relay")

// LaunchRelay sets up the connections to the Bitcoin and relay chains and
// launches the process of maintaining the relay.
func LaunchRelay(ctx context.Context) error {
	// TODO: Add connection to the Bitcoin chain:
	// btcChain, err := bitcoin.Connect(ctx, &maintainerConfig.Bitcoin)
	// if err != nil {
	// 	return fmt.Errorf("could not connect BTC chain: [%v]", err)
	// }

	// TODO: Add connection to the relay chain:
	// relayChain, err := connectRelayChain(config)
	// if err != nil {
	// 	return fmt.Errorf("could not connect relay chain: [%v]", err)
	// }

	newRelay(ctx, nil, nil)

	// TODO: Consider adding metrics.
	logger.Info("relay started")

	<-ctx.Done()
	return fmt.Errorf("unexpected context cancellation")
}

func newRelay(
	ctx context.Context,
	btcChain bitcoin.Chain,
	relayChain RelayChain,
) *Relay {
	relay := &Relay{
		btcChain:   btcChain,
		relayChain: relayChain,
	}

	go relay.startRelayControlLoop(ctx)

	return relay
}

// Relay is the part of maintainer responsible for maintaining the state of
// the relay on-chain contract.
type Relay struct {
	btcChain   bitcoin.Chain
	relayChain RelayChain
}

// startRelayControlLoop launches the loop responsible for controlling the relay.
func (r *Relay) startRelayControlLoop(ctx context.Context) {
	logger.Info("starting headers relay")

	defer func() {
		logger.Info("stopping headers relay")
	}()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			// TODO: Implement header relay loop. For now just print a message.
			logger.Info("relay is working")
		}

		time.Sleep(1 * time.Second)
	}
}
