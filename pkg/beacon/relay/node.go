package relay

import (
	"fmt"
	"math/big"
	"math/rand"
	"os"
	"reflect"
	"sync"

	"github.com/keep-network/keep-core/pkg/beacon/relay/candidate"
	relaychain "github.com/keep-network/keep-core/pkg/beacon/relay/chain"
	"github.com/keep-network/keep-core/pkg/beacon/relay/config"
	"github.com/keep-network/keep-core/pkg/beacon/relay/dkg"
	"github.com/keep-network/keep-core/pkg/beacon/relay/event"
	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/net"
	"github.com/keep-network/keep-core/pkg/thresholdgroup"
)

// Node represents the current state of a relay node.
type Node struct {
	mutex sync.Mutex

	// StakeID is the ID this node is using to prove its stake in the system.
	StakeID string
	// Staker is an on-chain identity that this node is using to prove its
	// stake in the system.
	Staker *candidate.Staker

	// External interactors.
	netProvider  net.Provider
	blockCounter chain.BlockCounter
	chainConfig  config.Chain

	// The IDs of the known stakes in the system, including this node's StakeID.
	stakeIDs      []string
	maxStakeIndex int

	groupPublicKeys [][]byte
	seenPublicKeys  map[string]struct{}
	myGroups        map[string]*membership
	pendingGroups   map[string]*membership
}

type membership struct {
	member  *thresholdgroup.Member
	channel net.BroadcastChannel
	index   int
}

// JoinGroupIfEligible takes a request id and the resulting entry value and
// checks if this client is one of the nodes elected by that entry to create a
// new relay group; if it is, this client enters the group creation process and,
// upon successfully completing it, submits the group public key to the passed
// relayChain. Note that this function returns immediately after determining
// whether the node is or is not eligible for the new group, and group joining
// and key submission is performed in a background goroutine.
func (n *Node) JoinGroupIfEligible(
	relayChain relaychain.Interface,
	requestID *big.Int,
	entryValue *big.Int,
) {
	if index := n.indexInEntryGroup(entryValue); index >= 0 {
		go func() {
			if !n.initializePendingGroup(requestID.String()) {
				// Failed to initialize; in progress for this entry.
				return
			}
			// Release control of this group if we error.
			defer n.flushPendingGroup(requestID.String())

			groupChannel, err := n.netProvider.ChannelFor(requestID.String())
			if err != nil {
				fmt.Fprintf(
					os.Stderr,
					"Error joining group channel for request group [%s]: [%v].\n",
					requestID.String(),
					err,
				)
				return
			}

			dkg.Init(groupChannel)
			member, err := dkg.ExecuteDKG(
				index,
				n.blockCounter,
				groupChannel,
				n.chainConfig.GroupSize,
				n.chainConfig.Threshold,
			)
			if err != nil {
				fmt.Fprintf(
					os.Stderr,
					"Failed DKG, error creating group: [%v].\n",
					err,
				)
				return
			}

			n.registerPendingGroup(requestID.String(), member, groupChannel)

			relayChain.SubmitGroupPublicKey(
				requestID,
				member.GroupPublicKeyBytes(),
			).OnComplete(func(registration *event.GroupRegistration, err error) {
				if err != nil {
					fmt.Fprintf(
						os.Stderr,
						"Failed submission of public key: [%v].\n",
						err,
					)
					return
				}

				n.RegisterGroup(
					registration.RequestID.String(),
					registration.GroupPublicKey,
				)
			})
		}()
	} else {
		err := n.SubmitTicketsForGroupSelection(entryValue.Bytes())
		if err != nil {
			fmt.Fprintf(
				os.Stderr,
				"Failed submission of tickets for group selection: [%v].\n",
				err,
			)
			return
		}
	}
}

// AddStaker registers a staker seen on-chain for the node's internal tracking.
func (n *Node) AddStaker(index int, staker string) {
	n.mutex.Lock()
	defer n.mutex.Unlock()

	if cap(n.stakeIDs) < index {
		// need something larger
		newSlice := make([]string, index*2)
		copy(newSlice, n.stakeIDs)
		n.stakeIDs = newSlice
	}

	n.stakeIDs[index] = staker

	if index > n.maxStakeIndex {
		n.maxStakeIndex = index
	}
}

// SyncStakingList performs an initial sync of the on-chain staker list into
// the node's internal state.
func (n *Node) SyncStakingList(stakingList []string) {
	for index, value := range stakingList {
		n.AddStaker(index, value)
	}
}

// RegisterGroup registers that a group was successfully created by the given
// requestID, and its group public key is groupPublicKey.
func (n *Node) RegisterGroup(requestID string, groupPublicKey []byte) {
	n.mutex.Lock()
	defer n.mutex.Unlock()

	// If we've already registered a group for this request ID, return early.
	if _, exists := n.seenPublicKeys[requestID]; exists {
		return
	}

	n.seenPublicKeys[requestID] = struct{}{}
	n.groupPublicKeys = append(n.groupPublicKeys, groupPublicKey)
	index := len(n.groupPublicKeys) - 1

	if membership, found := n.pendingGroups[requestID]; found && membership != nil {
		membership.index = index
		n.myGroups[requestID] = membership
		delete(n.pendingGroups, requestID)
	}
}

// initializePendingGroup grabs ownership of an attempt at group creation for a
// given goroutine. If it returns false, we're already in progress and failed to
// initialize.
func (n *Node) initializePendingGroup(requestID string) bool {
	n.mutex.Lock()
	defer n.mutex.Unlock()

	// If the pending group exists, we're already active
	if _, found := n.pendingGroups[requestID]; found {
		return false
	}

	// Pending group does not exist, take control
	n.pendingGroups[requestID] = nil

	return true
}

// flushPendingGroup if group creation fails, we clean our references to creating
// a group for a given request ID.
func (n *Node) flushPendingGroup(requestID string) {
	n.mutex.Lock()
	defer n.mutex.Unlock()

	if membership, found := n.pendingGroups[requestID]; found && membership == nil {
		delete(n.pendingGroups, requestID)
	}
}

// registerPendingGroup assigns a new membership for a given request ID.
// We overwrite our placeholder membership set by initializePendingGroup.
func (n *Node) registerPendingGroup(
	requestID string,
	member *thresholdgroup.Member,
	channel net.BroadcastChannel,
) {
	n.mutex.Lock()
	defer n.mutex.Unlock()

	if _, seen := n.seenPublicKeys[requestID]; seen {
		groupPublicKey := member.GroupPublicKeyBytes()
		// Start at the end since it's likely the public key was closer to the
		// end if it happened to come in before we had a chance to register it
		// as pending.
		existingIndex := len(n.groupPublicKeys) - 1
		for ; existingIndex >= 0; existingIndex-- {
			if reflect.DeepEqual(n.groupPublicKeys[existingIndex], groupPublicKey[:]) {
				break
			}
		}

		n.myGroups[requestID] = &membership{
			index:   existingIndex,
			member:  member,
			channel: channel,
		}
		delete(n.pendingGroups, requestID)
	} else {
		n.pendingGroups[requestID] = &membership{
			member:  member,
			channel: channel,
		}
	}
}

// Returns the 0-based index of this node in the group that will be spawned by
// the given entry value. If the node will not be in the group, returns -1.
func (n *Node) indexInEntryGroup(entryValue *big.Int) int {
	// FIXME By only using 64 bits, we're sacrificing a potentially large part
	// FIXME of the entry. We also need to reproduce this randomizer in
	// FIXME Solidity.
	shuffler := rand.New(rand.NewSource(entryValue.Int64()))

	n.mutex.Lock()
	shuffledStakeIDs := make([]string, n.maxStakeIndex+1)
	copy(shuffledStakeIDs, n.stakeIDs)
	currentStake := n.Staker.StakeID
	n.mutex.Unlock()

	shuffler.Shuffle(len(shuffledStakeIDs), func(i, j int) {
		shuffledStakeIDs[i], shuffledStakeIDs[j] = shuffledStakeIDs[j], shuffledStakeIDs[i]
	})

	for i, id := range shuffledStakeIDs {
		if id == currentStake {
			return i
		}
	}

	return -1
}
