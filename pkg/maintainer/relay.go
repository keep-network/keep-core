package maintainer

import (
	"context"
	"time"

	"github.com/ipfs/go-log"
	"github.com/keep-network/keep-core/pkg/bitcoin"
)

var logger = log.Logger("keep-maintainer-relay")

func newRelay(
	ctx context.Context,
	btcChain bitcoin.Chain,
	relayChain RelayChain,
) *Relay {
	relay := &Relay{
		btcChain:   btcChain,
		relayChain: relayChain,
	}

	go relay.startControlLoop(ctx)

	return relay
}

// Relay is the part of maintainer responsible for maintaining the state of
// the relay on-chain contract.
type Relay struct {
	btcChain   bitcoin.Chain
	relayChain RelayChain
}

// startControlLoop launches the loop responsible for controlling the relay.
func (r *Relay) startControlLoop(ctx context.Context) {
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
