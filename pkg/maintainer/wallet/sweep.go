package wallet

import (
	"context"
	"fmt"

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

	return wm.runOnceWalletUnlocked(
		ctx,
		walletPublicKeyHash,
		func() error {
			return coordinator.ProposeDepositsSweep(
				wm.chain,
				wm.btcChain,
				walletPublicKeyHash,
				0,
				deposits,
				true, // TODO: Change dry run argument to false
			)
		},
	)
}
