package tbtc

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"fmt"
	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/protocol/group"
	"github.com/keep-network/keep-core/pkg/tecdsa"
)

// WalletAction represents actions that can be performed by a wallet.
type WalletAction uint8

const (
	IdleWallet WalletAction = iota
	DepositSweep
	Redemption
	MovingFunds
	MovedFundsSweep
)

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
