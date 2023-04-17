package tbtc

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"fmt"
	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/protocol/group"
	"github.com/keep-network/keep-core/pkg/tecdsa"
	"golang.org/x/sync/semaphore"
	"math/big"
)

// WalletActionType represents actions types that can be performed by a wallet.
type WalletActionType uint8

const (
	Noop WalletActionType = iota
	DepositSweep
	Redemption
	MovingFunds
	MovedFundsSweep
)

// walletAction represents an action that can be performed by the wallet
// execution layer (i.e. the walletExecutor).
type walletAction interface {
	// run triggers the action execution. This function expects the signingExecutor
	// specific for the wallet executing the walletAction.
	run(signingExecutor *signingExecutor) error

	// actionType returns the specific type of the walletAction.
	actionType() WalletActionType
}

// errWalletExecutorBusy is an error returned when the walletExecutor
// cannot execute the submitted walletAction due to an ongoing work.
var errWalletExecutorBusy = fmt.Errorf("wallet executor is busy")

// walletExecutor is the execution layer for walletAction.
type walletExecutor struct {
	lock *semaphore.Weighted

	signingExecutor *signingExecutor
}

func newWalletExecutor(signingExecutor *signingExecutor) *walletExecutor {
	return &walletExecutor{
		lock:            semaphore.NewWeighted(1),
		signingExecutor: signingExecutor,
	}
}

// submit sends a walletAction to the walletExecutor. There is no guarantee
// the executor will actually perform the submitted action. The exact
// behavior depends on the executor internal state. In case the action
// is rejected by the executor, an error is returned.
func (we *walletExecutor) submit(action walletAction) error {
	if lockAcquired := we.lock.TryAcquire(1); !lockAcquired {
		return errWalletExecutorBusy
	}
	defer we.lock.Release(1)

	return action.run(we.signingExecutor)
}

// signBatch triggers signing of a messages batch. This method is basically a
// wrapper of the signingExecutor.signBatch method. It allows signing arbitrary
// data using the wallet execution layer, without the need of obtaining the
// lower-level signingExecutor. However, usage of this method should be limited
// to the minimum in favor of walletAction mechanism that should be the first
// choice, especially for more sophisticated actions.
func (we *walletExecutor) signBatch(
	ctx context.Context,
	messages []*big.Int,
	startBlock uint64,
) ([]*tecdsa.Signature, error) {
	return we.signingExecutor.signBatch(ctx, messages, startBlock)
}

// wallet represents a tBTC wallet. A wallet is one of the basic building
// blocks of the system that takes BTC under custody during the deposit
// process and gives that BTC back during redemptions.
type wallet struct {
	// publicKey is the unique ECDSA public key that identifies the
	// given wallet. This public key is also used to derive contract-specific
	// wallet identifiers (e.g. the Bridge contract identifies the wallet using
	// the SHA-256+RIPEMD-160 hash computed over the compressed ECDSA public key)
	publicKey *ecdsa.PublicKey
	// signingGroupOperators is the list holding operators' addresses that
	// form the whole wallet's signing group. This list may differ from the
	// original list outputted by the sortition protocol as it contains only
	// those signing group members who behaved properly during the DKG
	// protocol so all misbehaved members are not included here.
	// This list's size is always in the range [GroupQuorum, GroupSize].
	//
	// Each item in this list represents the given signing group member (seat)
	// and has a group.MemberIndex that is just the element's list index
	// incremented by one (e.g. element with index 0 has the group.MemberIndex
	// equal to 1 and so on).
	signingGroupOperators []chain.Address
}

// groupSize returns the actual size of the wallet's signing group. This
// value may be different from the GroupParameters.GroupSize parameter as some
// candidates may be excluded during distributed key generation.
func (w *wallet) groupSize() int {
	return len(w.signingGroupOperators)
}

// groupDishonestThreshold returns the dishonest threshold for the wallet's
// signing group. The returned value is computed using the wallet's actual
// signing group size for the given honest threshold provided as argument.
func (w *wallet) groupDishonestThreshold(honestThreshold int) int {
	return w.groupSize() - honestThreshold
}

func (w *wallet) String() string {
	publicKey := elliptic.Marshal(
		w.publicKey.Curve,
		w.publicKey.X,
		w.publicKey.Y,
	)

	return fmt.Sprintf(
		"wallet [0x%x] with a signing group of [%v]",
		publicKey,
		len(w.signingGroupOperators),
	)
}

// signer represents a threshold signer of a tBTC wallet. A signer holds
// a wallet tECDSA private key share and is able to participate in the
// signing process.
type signer struct {
	// wallet points to the tBTC wallet this signer belongs to.
	wallet wallet

	// signingGroupMemberIndex indicates the signer position (seat) in the
	// wallet signing group. Since the final wallet signing group may differ
	// from the original group outputted by the sortition protocol
	// (see wallet.signingGroupOperators documentation for reference), the
	// signingGroupMemberIndex may differ from the member index using
	// during the DKG protocol as well. The value of this index is in the
	// [1, len(wallet.signingGroupOperators)] range.
	signingGroupMemberIndex group.MemberIndex

	// privateKeyShare is the tECDSA private key share required to participate
	// in the signing process.
	privateKeyShare *tecdsa.PrivateKeyShare
}

// newSigner constructs a new instance of the wallet's signer.
func newSigner(
	walletPublicKey *ecdsa.PublicKey,
	walletSigningGroupOperators []chain.Address,
	signingGroupMemberIndex group.MemberIndex,
	privateKeyShare *tecdsa.PrivateKeyShare,
) *signer {
	wallet := wallet{
		publicKey:             walletPublicKey,
		signingGroupOperators: walletSigningGroupOperators,
	}

	return &signer{
		wallet:                  wallet,
		signingGroupMemberIndex: signingGroupMemberIndex,
		privateKeyShare:         privateKeyShare,
	}
}

func (s *signer) String() string {
	return fmt.Sprintf(
		"signer with index [%v] of %s",
		s.signingGroupMemberIndex,
		&s.wallet,
	)
}
