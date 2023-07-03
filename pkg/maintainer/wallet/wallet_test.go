package wallet

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/keep-network/keep-core/internal/testutils"
	"github.com/keep-network/keep-core/pkg/tbtc"
)

func TestRunIfWalletUnlocked_WhenLocked(t *testing.T) {
	localChain := NewLocalChain()

	walletPublicKeyHash := [20]byte{1}

	lockExpiration := time.Now().Add(500 * time.Millisecond)

	localChain.SetWalletLock(
		walletPublicKeyHash,
		lockExpiration,
		tbtc.ActionHeartbeat,
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
		tbtc.ActionDepositSweep,
		runFunc,
	)
	if err != nil {
		t.Fatal(err)
	}
}

func TestRunIfWalletUnlocked_WhenUnlocked(t *testing.T) {
	tests := map[string]struct{ expectedError error }{
		"no error in runFunc": {
			expectedError: nil,
		},
		"propagate error from runFunc": {
			expectedError: fmt.Errorf("boom, propagate up up up"),
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			localChain := NewLocalChain()

			walletPublicKeyHash := [20]byte{2}

			localChain.ResetWalletLock(walletPublicKeyHash)

			wasCalled := make(chan bool, 1)
			runFunc := func() error {
				wasCalled <- true
				return test.expectedError
			}

			wm := &walletMaintainer{
				config: Config{},
				chain:  localChain,
			}

			err := wm.runIfWalletUnlocked(
				context.Background(),
				walletPublicKeyHash,
				tbtc.ActionDepositSweep,
				runFunc,
			)

			<-wasCalled

			testutils.AssertErrorsSame(t, test.expectedError, err)
		})
	}
}
