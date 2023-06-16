package wallet

import (
	"context"
	"fmt"

	"time"

	"github.com/ipfs/go-log/v2"

	"github.com/keep-network/keep-core/internal/hexutils"
	"github.com/keep-network/keep-core/pkg/bitcoin"
)

var logger = log.Logger("keep-maintainer-wallet")

type walletMaintainer struct {
	config   Config
	chain    Chain
	btcChain bitcoin.Chain
}

// Initialize and start Wallet Coordination Maintainer.
func Initialize(
	parentCtx context.Context,
	config Config,
	chain Chain,
	btcChain bitcoin.Chain,

) {
	if config.RedemptionInterval == 0 {
		config.RedemptionInterval = DefaultRedemptionInterval
	}
	if config.SweepInterval == 0 {
		config.SweepInterval = DefaultSweepInterval
	}

	wm := &walletMaintainer{
		config:   config,
		chain:    chain,
		btcChain: btcChain,
	}

	go wm.startControlLoop(parentCtx)
}

// startControlLoop starts the loop responsible for controlling the wallet
// coordination maintainer.
func (wm *walletMaintainer) startControlLoop(ctx context.Context) {
	logger.Info("starting wallet coordination maintainer")
	defer logger.Info("stopping wallet coordination maintainer")

	initialRedemptionDelay := 5 * time.Second
	initialSweepDelay := 60 * time.Second

	redemptionTicker := time.NewTicker(initialRedemptionDelay)
	defer redemptionTicker.Stop()

	sweepTicker := time.NewTicker(initialSweepDelay)
	defer sweepTicker.Stop()

	logger.Infof("waiting [%s] until redemption task execution", initialRedemptionDelay)
	logger.Infof("waiting [%s] until sweep task execution", initialSweepDelay)

	for {
		select {
		case <-ctx.Done():
			logger.Info("tBTC Deposits Sweep Maintainer closed")
			return
		case <-redemptionTicker.C:
			logger.Info("starting redemption task execution...")

			// TODO: Implement

			logger.Infof("redemption task run completed; next run in [%s]", wm.config.RedemptionInterval)
			redemptionTicker.Reset(wm.config.RedemptionInterval)
		case <-sweepTicker.C:
			logger.Info("starting sweep task execution...")

			if err := wm.runSweepTask(ctx); err != nil {
				logger.Errorf("failed to run sweep task: [%v]", err)
			}

			logger.Infof("sweep task run completed; next run in [%s]", wm.config.SweepInterval)
			sweepTicker.Reset(wm.config.SweepInterval)
		}
	}
}

func (wm *walletMaintainer) runOnceWalletUnlocked(
	ctx context.Context,
	walletPublicKeyHash [20]byte,
	runFunc func() error,
) error {
	lockExpiration, walletAction, err := wm.chain.GetWalletLock(walletPublicKeyHash)
	if err != nil {
		return fmt.Errorf(
			"failed to get wallet lock for wallet public key hash [%s]: [%w]",
			hexutils.Encode(walletPublicKeyHash[:]),
			err,
		)
	}

	if lockExpiration.Unix() > 0 {
		logger.Infof(
			"wallet is locked due to [%s] action until [%s]; waiting...",
			walletAction.String(),
			lockExpiration.String(),
		)

		select {
		case <-ctx.Done():
			return nil
		case <-time.After(time.Until(lockExpiration)):
			break
		}
	}

	return runFunc()
}
