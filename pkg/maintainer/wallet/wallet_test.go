package wallet

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/keep-network/keep-core/pkg/tbtc"
)

func TestRunOnceWalletUnlocked_WhenLocked(t *testing.T) {
	localChain := newLocalChain()

	walletPublicKeyHash := [20]byte{1}

	lockExpiration := time.Now().Add(500 * time.Millisecond)

	localChain.setWalletLock(
		walletPublicKeyHash,
		lockExpiration,
		tbtc.Heartbeat,
	)

	funFunc := func() error {
		if time.Now().Before(lockExpiration) {
			return fmt.Errorf("too early")
		}

		return nil
	}

	wm := &walletMaintainer{
		config:   Config{},
		chain:    localChain,
		btcChain: newLocalBitcoinChain(),
	}

	err := wm.runOnceWalletUnlocked(
		context.Background(),
		walletPublicKeyHash,
		funFunc,
	)
	if err != nil {
		t.Fatal(err)
	}
}

func TestRunOnceWalletUnlocked_WhenUnlocked(t *testing.T) {
	localChain := newLocalChain()

	walletPublicKeyHash := [20]byte{2}

	lockExpiration := time.Now().Add(500 * time.Millisecond)

	localChain.resetWalletLock(walletPublicKeyHash)

	funFunc := func() error {
		if time.Now().After(lockExpiration) {
			return fmt.Errorf("too late")
		}

		return nil
	}

	wm := &walletMaintainer{
		config:   Config{},
		chain:    localChain,
		btcChain: newLocalBitcoinChain(),
	}

	err := wm.runOnceWalletUnlocked(
		context.Background(),
		walletPublicKeyHash,
		funFunc,
	)
	if err != nil {
		t.Fatal(err)
	}
}
