package signing

import (
	"github.com/ipfs/go-log"
	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/net"
	"github.com/keep-network/keep-core/pkg/protocol/group"
	"github.com/keep-network/keep-core/pkg/tecdsa"
	"math/big"
)

// Execute runs the tECDSA signing protocol, given a message to sign,
// broadcast channel to mediate with, a block counter used for time tracking,
// a member index to use in the group, private key share, dishonest threshold,
// and block height when signing protocol should start.
//
// This function also supports signing execution with a subset of the signing
// group by passing a non-empty excludedMembers slice holding the members that
// should be excluded.
func Execute(
	logger log.StandardLogger,
	message *big.Int,
	sessionID string,
	startBlockNumber uint64,
	memberIndex group.MemberIndex,
	privateKeyShare *tecdsa.PrivateKeyShare,
	groupSize int,
	dishonestThreshold int,
	excludedMembers []group.MemberIndex,
	blockCounter chain.BlockCounter,
	channel net.BroadcastChannel,
	membershipValidator *group.MembershipValidator,
) (*Result, error) {
	// TODO: Implementation.
	return &Result{
		R:          new(big.Int).Add(message, big.NewInt(1)),
		S:          new(big.Int).Add(message, big.NewInt(2)),
		RecoveryID: 1,
	}, nil
}