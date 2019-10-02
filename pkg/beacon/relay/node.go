package relay

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"math/big"
	"sync"

	relaychain "github.com/keep-network/keep-core/pkg/beacon/relay/chain"
	"github.com/keep-network/keep-core/pkg/beacon/relay/config"
	"github.com/keep-network/keep-core/pkg/beacon/relay/dkg"
	"github.com/keep-network/keep-core/pkg/beacon/relay/groupselection"
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
	chainConfig  *config.Chain

	groupRegistry *registry.Groups
}

// IsInGroup checks if this node is a member of the group which was selected to
// join a group which undergoes the process of generating a threshold relay entry.
func (n *Node) IsInGroup(groupPublicKey []byte) bool {
	return len(n.groupRegistry.GetGroup(groupPublicKey)) > 0
}

// JoinGroupIfEligible takes a threshold relay entry value and undergoes the
// process of joining a group if this node's virtual stakers prove eligible for
// the group generated by that entry. This is an interactive on-chain process,
// and JoinGroupIfEligible can block for an extended period of time while it
// completes the on-chain operation.
//
// Indirectly, the completion of the process is signaled by the formation of an
// on-chain group containing at least one of this node's virtual stakers.
func (n *Node) JoinGroupIfEligible(
	relayChain relaychain.Interface,
	signing chain.Signing,
	groupSelectionResult *groupselection.Result,
	newEntry *big.Int,
) {
	dkgStartBlockHeight := groupSelectionResult.GroupSelectionEndBlock

	for index, selectedStaker := range groupSelectionResult.SelectedStakers {
		// If we are amongst those chosen, kick off an instance of DKG. We may
		// have been selected multiple times (which would result in multiple
		// instances of DKG).
		if bytes.Compare(selectedStaker, n.Staker.ID()) == 0 {
			// capture player index for goroutine
			playerIndex := index

			// build the channel name and get the broadcast channel
			broadcastChannelName := channelNameForGroup(groupSelectionResult)

			// We should only join the broadcast channel if we're
			// elligible for the group
			broadcastChannel, err := n.netProvider.ChannelFor(
				broadcastChannelName,
			)
			if err != nil {
				logger.Errorf(
					"failed to get broadcastChannel for name [%s] with err: [%v]",
					broadcastChannelName,
					err,
				)
				return
			}

			go func() {
				signer, err := dkg.ExecuteDKG(
					newEntry,
					playerIndex,
					n.chainConfig.GroupSize,
					n.chainConfig.DishonestThreshold(),
					dkgStartBlockHeight,
					n.blockCounter,
					relayChain,
					signing,
					broadcastChannel,
				)
				if err != nil {
					logger.Errorf("failed to execute dkg: [%v]", err)
					return
				}

				err = n.groupRegistry.RegisterGroup(
					signer,
					broadcastChannelName,
				)
				if err != nil {
					logger.Errorf("failed to register a group: [%v]", err)
				}
			}()
		}
	}

	return
}

// channelNameForGroup takes the selected stakers, and does the
// following to construct the broadcastChannel name:
// * concatenates all of the staker values
// * returns the hashed concatenated values in hexadecimal representation
func channelNameForGroup(group *groupselection.Result) string {
	var channelNameBytes []byte
	for _, staker := range group.SelectedStakers {
		channelNameBytes = append(channelNameBytes, staker...)
	}

	hash := sha256.Sum256(channelNameBytes)
	hexChannelName := hex.EncodeToString(hash[:])

	return hexChannelName
}
