package maintainer

import (
	"context"
	"time"

	"github.com/ipfs/go-log"
	"github.com/keep-network/keep-core/pkg/bitcoin"
)

var logger = log.Logger("keep-maintainer-bitcoin-difficulty")

func newBitcoinDifficultyMaintainer(
	ctx context.Context,
	btcChain bitcoin.Chain,
	chain BitcoinDifficultyChain,
) *BitcoinDifficultyMaintainer {
	bitcoinDifficultyMaintainer := &BitcoinDifficultyMaintainer{
		btcChain: btcChain,
		chain:    chain,
	}

	go bitcoinDifficultyMaintainer.startControlLoop(ctx)

	return bitcoinDifficultyMaintainer
}

// BitcoinDifficultyMaintainer is the part of maintainer responsible for
// maintaining the state of the Bitcoin difficulty on-chain contract.
type BitcoinDifficultyMaintainer struct {
	btcChain bitcoin.Chain
	chain    BitcoinDifficultyChain
}

// startControlLoop launches the loop responsible for controlling the Bitcoin
// difficulty maintainer.
func (r *BitcoinDifficultyMaintainer) startControlLoop(ctx context.Context) {
	logger.Info("starting Bitcoin difficulty maintainer")

	defer func() {
		logger.Info("stopping Bitcoin difficulty maintainer")
	}()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			// TODO: Implement the maintainer loop. For now just print a message.
			logger.Info("Bitcoin difficulty maintainer is working")
		}

		time.Sleep(1 * time.Second)
	}
}
