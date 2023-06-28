package wallet

import (
	"context"
	"fmt"

	"time"

	"github.com/ipfs/go-log/v2"

	"github.com/keep-network/keep-core/internal/hexutils"
	"github.com/keep-network/keep-core/pkg/bitcoin"
	"github.com/keep-network/keep-core/pkg/tbtc"
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
	if config.DepositSweepInterval == 0 {
		config.DepositSweepInterval = DefaultDepositSweepInterval
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
	initialDepositSweepDelay := 60 * time.Second

	redemptionTicker := time.NewTicker(initialRedemptionDelay)
	defer redemptionTicker.Stop()

	depositSweepTicker := time.NewTicker(initialDepositSweepDelay)
	defer depositSweepTicker.Stop()

	logger.Infof("waiting [%s] until redemption task execution", initialRedemptionDelay)
	logger.Infof("waiting [%s] until deposit sweep task execution", initialDepositSweepDelay)

	for {
		select {
		case <-ctx.Done():
			return
		case <-redemptionTicker.C:
			// Set the ticker to the expected interval.
			redemptionTicker.Reset(wm.config.RedemptionInterval)

			logger.Info("starting redemption task execution...")

			// TODO: Implement

			logger.Infof(
				"redemption task run completed; next run in [%s]",
				wm.config.RedemptionInterval,
			)
		case <-depositSweepTicker.C:
			// Set the ticker to the expected interval.
			depositSweepTicker.Reset(wm.config.DepositSweepInterval)

			logger.Info("starting deposit sweep task execution...")

			if err := wm.runDepositSweepTask(ctx); err != nil {
				logger.Errorf("failed to run deposit sweep task: [%v]", err)
			}

			logger.Infof(
				"deposit sweep task run completed; next run in [%s]",
				wm.config.DepositSweepInterval,
			)
		}
	}
}

func (wm *walletMaintainer) runIfWalletUnlocked(
	ctx context.Context,
	walletPublicKeyHash [20]byte,
	currentWalletAction tbtc.WalletActionType,
	runFunc func() error,
) error {
	lockExpiration, lockWalletAction, err := wm.chain.GetWalletLock(walletPublicKeyHash)
	if err != nil {
		return fmt.Errorf(
			"failed to get wallet lock for wallet public key hash [%s]: [%w]",
			hexutils.Encode(walletPublicKeyHash[:]),
			err,
		)
	}

	if lockExpiration.After(time.Now()) {
		logger.Infof(
			"wallet is locked due to [%s] action until [%s]; skipping [%s] execution...",
			lockWalletAction.String(),
			lockExpiration.String(),
			currentWalletAction.String(),
		)
		return nil
	}

	return runFunc()
}
