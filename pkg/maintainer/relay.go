package maintainer

import (
	"context"
	"time"

	"github.com/ipfs/go-log"
	"github.com/keep-network/keep-core/pkg/bitcoin"
)

var logger = log.Logger("keep-maintainer-relay")

// RelayChain is an interface that provides the ability to communicate with the
// relay on-chain contract.
type RelayChain interface {
	// Retarget adds a new epoch to the relay by providing a proof
	// of the difficulty before and after the retarget.
	Retarget(headers []bitcoin.BlockHeader) error
}

func newRelay(
	ctx context.Context,
	btcChain bitcoin.Chain,
	chain RelayChain,
) *Relay {
	relay := &Relay{
		btcChain: btcChain,
		chain:    chain,
	}

	go relay.startControlLoop(ctx)

	return relay
}

// Relay is the part of maintainer responsible for maintaining the state of
// the relay on-chain contract.
type Relay struct {
	btcChain bitcoin.Chain
	chain    RelayChain
}

// startControlLoop launches the loop responsible for controlling the relay.
func (r *Relay) startControlLoop(ctx context.Context) {
	logger.Info("starting Bitcoin difficulty maintainer")

	defer func() {
		logger.Info("stopping Bitcoin difficulty maintainer")
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
