package wallet

import (
	"context"
	"fmt"
	"time"

	"github.com/keep-network/keep-core/internal/hexutils"
	"github.com/keep-network/keep-core/pkg/coordinator"
)

func (wm *walletMaintainer) runSweepTask(ctx context.Context) error {
	depositSweepMaxSize, err := wm.chain.GetDepositSweepMaxSize()
	if err != nil {
		return fmt.Errorf("failed to get deposit sweep max size: [%v]", err)
	}

	walletPublicKeyHash, deposits, err := coordinator.FindDepositsToSweep(
		wm.chain,
		wm.btcChain,
		[20]byte{},
		depositSweepMaxSize,
	)
	if err != nil {
		return fmt.Errorf("failed to prepare deposits sweep proposal: %v", err)
	}

	lockExpiration, walletAction, err := wm.chain.GetWalletLock(walletPublicKeyHash)
	if err != nil {
		return fmt.Errorf(
			"failed to get wallet locks for wallet public key hash[%s]: [%w]",
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

	return coordinator.ProposeDepositsSweep(
		wm.chain,
		wm.btcChain,
		walletPublicKeyHash,
		0,
		deposits,
		true, // TODO: Change dry run argument to false
	)
}
