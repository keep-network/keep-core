package inactivity

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"sort"

	"github.com/ipfs/go-log/v2"

	"github.com/keep-network/keep-core/pkg/net"
	"github.com/keep-network/keep-core/pkg/protocol/group"
	"github.com/keep-network/keep-core/pkg/protocol/state"
)

// ClaimPreimage represents an inactivity claim preimage.
type ClaimPreimage struct {
	Nonce                  *big.Int
	WalletPublicKey        *ecdsa.PublicKey
	InactiveMembersIndexes []group.MemberIndex
	HeartbeatFailed        bool
}

func NewClaimPreimage(
	nonce *big.Int,
	walletPublicKey *ecdsa.PublicKey,
	inactiveMembersIndexes []group.MemberIndex,
	heartbeatFailed bool,
) *ClaimPreimage {
	// Made the inactive member indexes unique as expected by the on-chain
	// contract.
	indexesCache := make(map[group.MemberIndex]bool)
	uniqueIndexes := []group.MemberIndex{}

	for _, index := range inactiveMembersIndexes {
		if _, exists := indexesCache[index]; !exists {
			indexesCache[index] = true
			uniqueIndexes = append(uniqueIndexes, index)
		}
	}

	// Sort the inactive member indexes as expected by the on-chain contract.
	sort.Slice(uniqueIndexes, func(i, j int) bool {
		return uniqueIndexes[i] < uniqueIndexes[j]
	})

	return &ClaimPreimage{
		Nonce:                  nonce,
		WalletPublicKey:        walletPublicKey,
		InactiveMembersIndexes: uniqueIndexes,
		HeartbeatFailed:        heartbeatFailed,
	}
}

const ClaimHashByteSize = 32

// ClaimHash is a hash of the inactivity claim. The hashing algorithm used
// depends on the client code.
type ClaimHash [ClaimHashByteSize]byte

// ClaimHashFromBytes converts bytes slice to ClaimHash. It requires provided
// bytes slice size to be exactly ClaimHashByteSize.
func ClaimHashFromBytes(bytes []byte) (ClaimHash, error) {
	var hash ClaimHash

	if len(bytes) != ClaimHashByteSize {
		return hash, fmt.Errorf(
			"bytes length is not equal %v", ClaimHashByteSize,
		)
	}
	copy(hash[:], bytes)

	return hash, nil
}

// SignedClaimHash represents information pertaining to the process of signing
// an inactivity claim: the public key used during signing, the resulting
// signature and the hash of the inactivity claim that was used during signing.
type SignedClaimHash struct {
	PublicKey []byte
	Signature []byte
	ClaimHash ClaimHash
}

type ClaimSigner interface {
	SignClaim(claim *ClaimPreimage) (*SignedClaimHash, error)
	VerifySignature(signedClaim *SignedClaimHash) (bool, error)
}

type ClaimSubmitter interface {
	SubmitClaim(
		ctx context.Context,
		memberIndex group.MemberIndex,
		claim *ClaimPreimage,
		signatures map[group.MemberIndex][]byte,
	) error
}

func PublishClaim(
	ctx context.Context,
	logger log.StandardLogger,
	sessionID string,
	memberIndex group.MemberIndex,
	channel net.BroadcastChannel,
	groupSize int,
	dishonestThreshold int,
	membershipValidator *group.MembershipValidator,
	claimSigner ClaimSigner,
	claimSubmitter ClaimSubmitter,
	claim *ClaimPreimage,
) error {
	initialState := &claimSigningState{
		BaseAsyncState: state.NewBaseAsyncState(),
		channel:        channel,
		claimSigner:    claimSigner,
		claimSubmitter: claimSubmitter,
		member: newSigningMember(
			logger,
			memberIndex,
			groupSize,
			dishonestThreshold,
			membershipValidator,
			sessionID,
		),
		claim: claim,
	}

	stateMachine := state.NewAsyncMachine(logger, ctx, channel, initialState)

	lastState, err := stateMachine.Execute()
	if err != nil {
		return err
	}

	_, ok := lastState.(*claimSubmissionState)
	if !ok {
		return fmt.Errorf("execution ended on state %T", lastState)
	}

	return nil
}

// RegisterUnmarshallers initializes the given broadcast channel to be able to
// perform inactivity claim interactions by registering all the required
// protocol message unmarshallers.
func RegisterUnmarshallers(channel net.BroadcastChannel) {
	channel.SetUnmarshaler(func() net.TaggedUnmarshaler {
		return &claimSignatureMessage{}
	})
}
