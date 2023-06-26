package spv

import (
	"context"
	"time"

	"github.com/ipfs/go-log/v2"

	"github.com/keep-network/keep-core/pkg/bitcoin"
)

var logger = log.Logger("keep-maintainer-spv")

func Initialize(
	ctx context.Context,
	config Config,
	chain Chain,
	btcChain bitcoin.Chain,
) {
	spvMaintainer := &spvMaintainer{
		config:   config,
		chain:    chain,
		btcChain: btcChain,
	}

	go spvMaintainer.startControlLoop(ctx)
}

type spvMaintainer struct {
	config   Config
	chain    Chain
	btcChain bitcoin.Chain
}

func (sm *spvMaintainer) startControlLoop(ctx context.Context) {
	logger.Info("starting SPV maintainer")

	defer func() {
		logger.Info("stopping SPV maintainer")
	}()

	for {
		err := sm.maintainSpv(ctx)
		if err != nil {
			logger.Errorf(
				"error while maintaining SPV: [%v]; restarting maintainer",
				err,
			)
		}

		select {
		case <-time.After(sm.config.RestartBackOffTime):
		case <-ctx.Done():
			return
		}
	}
}

func (sm *spvMaintainer) maintainSpv(ctx context.Context) error {
	logger.Infof("Maintaining SPV proof...")

	// TODO: Implement. For now, just wait.
	<-ctx.Done()

	return nil
}
