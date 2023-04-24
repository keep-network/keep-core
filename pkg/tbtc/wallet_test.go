package tbtc

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"github.com/keep-network/keep-core/pkg/internal/testutils"
	"github.com/keep-network/keep-core/pkg/tecdsa"
	"math/big"
	"testing"
	"time"
)

func TestWalletDispatcher_Dispatch(t *testing.T) {
	walletDispatcher := newWalletDispatcher()

	wallet1 := generateWallet(big.NewInt(100))
	wallet2 := generateWallet(big.NewInt(101))

	// Ctx for first actions of both wallets.
	ctxActions1, cancelCtxActions1 := context.WithCancel(context.Background())
	defer cancelCtxActions1()
	// Ctx for second actions of both wallets.
	ctxActions2, cancelCtxActions2 := context.WithCancel(context.Background())
	defer cancelCtxActions2()

	wallet1Action1 := &mockWalletAction{
		executeFn: func() error {
			<-ctxActions1.Done()
			return nil // complete with success
		},
		actionWallet: wallet1,
	}
	wallet1Action2 := &mockWalletAction{
		executeFn: func() error {
			<-ctxActions2.Done()
			return nil // complete with success
		},
		actionWallet: wallet1,
	}
	wallet2Action1 := &mockWalletAction{
		executeFn: func() error {
			<-ctxActions1.Done()
			return fmt.Errorf("unexpected error") // complete with error
		},
		actionWallet: wallet2,
	}
	wallet2Action2 := &mockWalletAction{
		executeFn: func() error {
			<-ctxActions2.Done()
			return nil // complete with success
		},
		actionWallet: wallet2,
	}

	// Dispatch Action 1 for Wallet 1.
	err := walletDispatcher.dispatch(wallet1Action1)
	if err != nil {
		t.Errorf("unexpected error: [%v]", err)
	}

	// Another Action 1 for Wallet 2.
	err = walletDispatcher.dispatch(wallet2Action1)
	if err != nil {
		t.Errorf("unexpected error: [%v]", err)
	}

	// Try to dispatch Action 1 for Wallet 1 again.
	err = walletDispatcher.dispatch(wallet1Action1)
	testutils.AssertErrorsSame(t, errWalletBusy, err)

	// Try to dispatch Action 1 for Wallet 2 again.
	err = walletDispatcher.dispatch(wallet2Action1)
	testutils.AssertErrorsSame(t, errWalletBusy, err)

	// Try to dispatch Action 2 for Wallet 1.
	err = walletDispatcher.dispatch(wallet1Action2)
	testutils.AssertErrorsSame(t, errWalletBusy, err)

	// Try to dispatch Action 2 for Wallet 2.
	err = walletDispatcher.dispatch(wallet2Action2)
	testutils.AssertErrorsSame(t, errWalletBusy, err)

	// Complete dispatched actions.
	cancelCtxActions1()
	<-ctxActions1.Done()

	// Give some time to release the lock.
	time.Sleep(1 * time.Second)

	// Dispatch Action 2 for Wallet 1.
	err = walletDispatcher.dispatch(wallet1Action2)
	if err != nil {
		t.Errorf("unexpected error: [%v]", err)
	}

	// Dispatch Action 2 for Wallet 2.
	err = walletDispatcher.dispatch(wallet2Action2)
	if err != nil {
		t.Errorf("unexpected error: [%v]", err)
	}
}

type mockWalletAction struct {
	executeFn    func() error
	actionWallet wallet
}

func (mwa *mockWalletAction) execute() error {
	return mwa.executeFn()
}

func (mwa *mockWalletAction) wallet() wallet {
	return mwa.actionWallet
}

func (mwa *mockWalletAction) actionType() WalletActionType {
	return Noop
}

func generateWallet(privateKey *big.Int) wallet {
	x, y := tecdsa.Curve.ScalarBaseMult(privateKey.Bytes())
	publicKey := &ecdsa.PublicKey{
		Curve: tecdsa.Curve,
		X:     x,
		Y:     y,
	}

	return wallet{
		publicKey: publicKey,
	}
}
