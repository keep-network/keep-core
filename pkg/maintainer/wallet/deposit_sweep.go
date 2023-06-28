package wallet

import (
	"context"
	"fmt"

	"github.com/keep-network/keep-core/pkg/coordinator"
	"github.com/keep-network/keep-core/pkg/tbtc"
)

func (wm *walletMaintainer) runDepositSweepTask(ctx context.Context) error {
	depositSweepMaxSize, err := wm.chain.GetDepositSweepMaxSize()
	if err != nil {
		return fmt.Errorf("failed to get deposit sweep max size: [%w]", err)
	}

	walletPublicKeyHash, deposits, err := coordinator.FindDepositsToSweep(
		wm.chain,
		wm.btcChain,
		[20]byte{},
		depositSweepMaxSize,
	)
	if err != nil {
		return fmt.Errorf("failed to prepare deposits sweep proposal: [%w]", err)
	}

	if len(deposits) == 0 {
		logger.Info("no deposits to sweep")
		return nil
	}

	return wm.runIfWalletUnlocked(
		ctx,
		walletPublicKeyHash,
		tbtc.DepositSweep,
		func() error {
			return coordinator.ProposeDepositsSweep(
				wm.chain,
				wm.btcChain,
				walletPublicKeyHash,
				0,
				deposits,
				false,
			)
		},
	)
}
