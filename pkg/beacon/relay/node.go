package relay

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"math/big"
	"sync"

	bn256 "github.com/ethereum/go-ethereum/crypto/bn256/cloudflare"
	"github.com/keep-network/keep-core/pkg/altbn128"

	relayChain "github.com/keep-network/keep-core/pkg/beacon/relay/chain"
	"github.com/keep-network/keep-core/pkg/beacon/relay/group"

	relaychain "github.com/keep-network/keep-core/pkg/beacon/relay/chain"
	"github.com/keep-network/keep-core/pkg/beacon/relay/dkg"
	"github.com/keep-network/keep-core/pkg/beacon/relay/registry"
	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/net"
)

// Node represents the current state of a relay node.
type Node struct {
	mutex sync.Mutex

	// Staker is an on-chain identity that this node is using to prove its
	// stake in the system.
	Staker chain.Staker

	// External interactors.
	netProvider  net.Provider
	blockCounter chain.BlockCounter
	chainConfig  *relaychain.Config

	groupRegistry *registry.Groups
}

// IsInGroup checks if this node is a member of the group which was selected to
// join a group which undergoes the process of generating a threshold relay entry.
func (n *Node) IsInGroup(groupPublicKey []byte) bool {
	return len(n.groupRegistry.GetGroup(groupPublicKey)) > 0
}

// JoinDKGIfEligible takes a seed value and undergoes the process of the
// distributed key generation if this node's operator proves eligible for
// the group generated by that seed. This is an interactive on-chain process,
// and JoinDKGIfEligible can block for an extended period of time while it
// completes the on-chain operation.
func (n *Node) JoinDKGIfEligible(
	relayChain relaychain.Interface,
	signing chain.Signing,
	dkgSeed *big.Int,
	dkgStartBlockNumber uint64,
) {
	// TODO: Fetch selected group members from the sortition pool.
	groupMembers := make([]relaychain.StakerAddress, 0)

	if len(groupMembers) > maxGroupSize {
		logger.Errorf(
			"group size larger than supported: [%v]",
			len(groupMembers),
		)
		return
	}

	indexes := make([]uint8, 0)
	for index, groupMember := range groupMembers {
		// See if we are amongst those chosen
		if bytes.Compare(groupMember, n.Staker.Address()) == 0 {
			indexes = append(indexes, uint8(index))
		}
	}

	// create temporary broadcast channel name for DKG using the
	// group selection seed
	channelName := dkgSeed.Text(16)

	if len(indexes) > 0 {
		broadcastChannel, err := n.netProvider.BroadcastChannelFor(channelName)
		if err != nil {
			logger.Errorf("failed to get broadcast channel: [%v]", err)
			return
		}

		// TODO: Setup the correct validator.
		membershipValidator := group.NewStakersMembershipValidator(
			nil,
			signing,
		)

		err = broadcastChannel.SetFilter(membershipValidator.IsInGroup)
		if err != nil {
			logger.Errorf(
				"could not set filter for channel [%v]: [%v]",
				broadcastChannel.Name(),
				err,
			)
		}

		for _, index := range indexes {
			// Capture the player index for the goroutine.
			playerIndex := index

			go func() {
				signer, err := dkg.ExecuteDKG(
					dkgSeed,
					playerIndex,
					n.chainConfig.GroupSize,
					n.chainConfig.DishonestThreshold(),
					membershipValidator,
					dkgStartBlockNumber,
					n.blockCounter,
					relayChain,
					signing,
					broadcastChannel,
				)
				if err != nil {
					logger.Errorf("failed to execute dkg: [%v]", err)
					return
				}

				// The final broadcast channel name for group is the compressed
				// public key of the group.
				channelName := hex.EncodeToString(
					signer.GroupPublicKeyBytesCompressed(),
				)

				// Register the candidate group. Note that such a group is
				// non-operable and the node should monitor for the group
				// approval event in order to register an approved group in the
				// group registry as well.
				err = n.groupRegistry.RegisterCandidateGroup(signer, channelName)
				if err != nil {
					logger.Errorf(
						"[member:%v] failed to register a candidate group [%v]: [%v]",
						signer.MemberID(),
						channelName,
						err,
					)
				}

				logger.Infof(
					"[member:%v] candidate group [%v] registered successfully",
					signer.MemberID(),
					channelName,
				)
			}()
		}
	}

	return
}

// ForwardSignatureShares enables the ability to forward signature shares
// messages to other nodes even if this node is not a part of the group which
// signs the relay entry.
func (n *Node) ForwardSignatureShares(groupPublicKeyBytes []byte) {
	name, err := channelNameForPublicKeyBytes(groupPublicKeyBytes)
	if err != nil {
		logger.Warningf("could not forward signature shares: [%v]", err)
		return
	}

	n.netProvider.BroadcastChannelForwarderFor(name)
}

// ResumeSigningIfEligible enables a client to rejoin the ongoing signing process
// after it was crashed or restarted and if it belongs to the signing group.
func (n *Node) ResumeSigningIfEligible(
	relayChain relayChain.Interface,
	signing chain.Signing,
) {
	isEntryInProgress, err := relayChain.IsEntryInProgress()
	if err != nil {
		logger.Errorf(
			"failed checking if an entry is in progress: [%v]",
			err,
		)
		return
	}

	if isEntryInProgress {
		previousEntry, err := relayChain.CurrentRequestPreviousEntry()
		if err != nil {
			logger.Errorf(
				"failed to get a previous entry for the current request: [%v]",
				err,
			)
			return
		}
		entryStartBlock, err := relayChain.CurrentRequestStartBlock()
		if err != nil {
			logger.Errorf(
				"failed to get a start block for the current request: [%v]",
				err,
			)
			return
		}
		groupPublicKey, err := relayChain.CurrentRequestGroupPublicKey()
		if err != nil {
			logger.Errorf(
				"failed to get a group public key for the current request: [%v]",
				err,
			)
			return
		}

		logger.Infof(
			"attempting to rejoin the current signing process [0x%x]",
			groupPublicKey,
		)
		n.GenerateRelayEntry(
			previousEntry,
			relayChain,
			signing,
			groupPublicKey,
			entryStartBlock.Uint64(),
		)
	}
}

// channelNameForPublicKey takes group public key represented by marshalled
// G2 point and transforms it into a broadcast channel name.
// Broadcast channel name for group is the hexadecimal representation of
// compressed public key of the group.
func channelNameForPublicKeyBytes(groupPublicKey []byte) (string, error) {
	g2 := new(bn256.G2)

	if _, err := g2.Unmarshal(groupPublicKey); err != nil {
		return "", fmt.Errorf("could not create channel name: [%v]", err)
	}

	return channelNameForPublicKey(g2), nil
}

// channelNameForPublicKey takes group public key represented by G2 point
// and transforms it into a broadcast channel name.
// Broadcast channel name for group is the hexadecimal representation of
// compressed public key of the group.
func channelNameForPublicKey(groupPublicKey *bn256.G2) string {
	altbn128GroupPublicKey := altbn128.G2Point{G2: groupPublicKey}
	return hex.EncodeToString(altbn128GroupPublicKey.Compress())
}
