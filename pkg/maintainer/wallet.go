package maintainer

import (
	"context"
	"fmt"
	"time"

	"github.com/ipfs/go-log"
	"github.com/keep-network/keep-core/pkg/bitcoin"
)

const (
	// Default value for back-off time which should be applied when the wallet
	// maintainer is restarted. It helps to avoid being flooded with error logs
	//  in case of a permanent error in the wallet maintainer.
	defaultRestartBackoffTime = 120 * time.Second
)

// TODO: Store the logger inside the wallet maintainer.
var walletLogger = log.Logger("keep-maintainer-wallet")

func newWalletMaintainer(
	ctx context.Context,
	chain WalletChain,
	restartBackOffTime time.Duration,
) *WalletMaintainer {
	walletMaintainer := &WalletMaintainer{
		chain:              chain,
		restartBackOffTime: restartBackOffTime,
	}

	go walletMaintainer.startControlLoop(ctx)

	return walletMaintainer
}

// TODO: Description
type WalletMaintainer struct {
	chain WalletChain

	restartBackOffTime time.Duration
}

// startControlLoop launches the loop responsible for controlling the wallet
// maintainer.
func (wm *WalletMaintainer) startControlLoop(ctx context.Context) {
	walletLogger.Info("starting wallet maintainer")

	defer func() {
		walletLogger.Info("stopping wallet maintainer")
	}()

	for {
		err := wm.checkWallet(ctx)
		if err != nil {
			walletLogger.Errorf(
				"restarting wallet maintainer due to error while checking "+
					"wallet: [%v]",
				err,
			)
		}

		select {
		case <-time.After(wm.restartBackOffTime):
		case <-ctx.Done():
			return
		}
	}
}

// TODO: Description
func (wm *WalletMaintainer) checkWallet(ctx context.Context) error {
	eligible, err := wm.checkRequestNewWalletEligibility()
	if err != nil {
		return fmt.Errorf(
			"failed to check new wallet request eligibility: [%w]",
			err,
		)
	}

	if !eligible {
		return nil
	}

	if err := wm.requestNewWallet(); err != nil {
		return fmt.Errorf(
			"failed to request a new wallet: [%w]",
			err,
		)
	}

	return nil
}

// TODO: Description
func (wm *WalletMaintainer) checkRequestNewWalletEligibility() (bool, error) {
	state, err := wm.chain.GetWalletCreationState()
	if err != nil {
		return false, fmt.Errorf(
			"cannot get wallet creation state: [%w]", err,
		)
	}

	if state != Idle {
		return false, nil
	}

	activeWalletPubKeyHash, err := wm.chain.ActiveWalletPubKeyHash()
	if err != nil {
		return false, fmt.Errorf(
			"cannot get active wallet public key hash: [%w]", err,
		)
	}

	if activeWalletPubKeyHash == [20]byte{} {
		return true, nil
	}

	walletPublicKey, mainUtxoHash, createdAt, err := wm.chain.GetWalletInfo(
		activeWalletPubKeyHash,
	)
	if err != nil {
		return false, fmt.Errorf("cannot get active wallet info: [%w]", err)
	}

	walletCreationPeriod,
		walletCreationMinBtcBalance,
		walletCreationMaxBtcBalance,
		err := wm.chain.WalletParameters()
	if err != nil {
		return false, fmt.Errorf("cannot get wallet creation period: [%w]", err)
	}

	walletMainUtxo, err := wm.getWalletMainUtxo(mainUtxoHash, walletPublicKey)
	if err != nil {
		return false, fmt.Errorf("cannot get wallet main UTXO: [%w]", err)
	}

	activeWalletBalance := uint64(walletMainUtxo.Value)

	activeWalletOldEnough := getUnixTime() > createdAt+walletCreationPeriod

	if activeWalletOldEnough &&
		activeWalletBalance >= walletCreationMinBtcBalance {
		return true, nil
	}

	if activeWalletBalance >= walletCreationMaxBtcBalance {
		return true, nil
	}

	return false, nil
}

// TODO: Description
func (wm *WalletMaintainer) requestNewWallet() error {
	// TODO: Implement
	return nil
}

// TODO: Description
func (wm *WalletMaintainer) getWalletMainUtxo(
	mainUtxoHash [32]byte,
	walletPublicKey []byte,
) (bitcoin.UnspentTransactionOutput, error) {
	// TODO: Connect to the Bitcoin chain, retrieve all the UTXOs of
	//       `walletPublicKey`, find which one of them matches `mainUtxoHash`
	//       and return it.

	return bitcoin.UnspentTransactionOutput{}, nil
}

func getUnixTime() uint32 {
	// TODO: Consider replacing it with a real Ethereum block timestamp.
	return uint32(time.Now().Unix())
}
