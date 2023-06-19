package wallet

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/keep-network/keep-core/pkg/tbtc"
)

func TestRunIfWalletUnlocked_WhenLocked(t *testing.T) {
	localChain := newLocalChain()

	walletPublicKeyHash := [20]byte{1}

	lockExpiration := time.Now().Add(500 * time.Millisecond)

	localChain.setWalletLock(
		walletPublicKeyHash,
		lockExpiration,
		tbtc.Heartbeat,
	)

	runFunc := func() error {
		return fmt.Errorf("boom, you should not run me")
	}

	wm := &walletMaintainer{
		config: Config{},
		chain:  localChain,
	}

	err := wm.runIfWalletUnlocked(
		context.Background(),
		walletPublicKeyHash,
		tbtc.DepositSweep,
		runFunc,
	)
	if err != nil {
		t.Fatal(err)
	}
}

func TestRunIfWalletUnlocked_WhenUnlocked(t *testing.T) {
	localChain := newLocalChain()

	walletPublicKeyHash := [20]byte{2}

	localChain.resetWalletLock(walletPublicKeyHash)

	runFunc := func() error {
		return nil
	}

	wm := &walletMaintainer{
		config: Config{},
		chain:  localChain,
	}

	err := wm.runIfWalletUnlocked(
		context.Background(),
		walletPublicKeyHash,
		tbtc.DepositSweep,
		runFunc,
	)
	if err != nil {
		t.Fatal(err)
	}
}
