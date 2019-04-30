package relay

import (
	"fmt"
	"os"
	"sync"

	relaychain "github.com/keep-network/keep-core/pkg/beacon/relay/chain"
	"github.com/keep-network/keep-core/pkg/beacon/relay/dkg"
	"github.com/keep-network/keep-core/pkg/net"
)

// GroupRegistry represents a collection of Keep groups.
type GroupRegistry struct {
	mutex sync.Mutex

	myGroups map[string][]*membership

	relayChain relaychain.Interface
}

type membership struct {
	signer  *dkg.ThresholdSigner
	channel net.BroadcastChannel
}

// NewGroupRegistry returns an empty GroupRegistry.
func NewGroupRegistry(
	relayChain relaychain.Interface,
) GroupRegistry {
	return GroupRegistry{
		myGroups:   make(map[string][]*membership),
		relayChain: relayChain,
	}
}

// RegisterGroup registers that a group was successfully created by the given
// groupPublicKey.
func (gr *GroupRegistry) RegisterGroup(signer *dkg.ThresholdSigner,
	channel net.BroadcastChannel) {

	gr.mutex.Lock()
	defer gr.mutex.Unlock()

	groupPublicKey := string(signer.GroupPublicKeyBytes())

	gr.myGroups[groupPublicKey] = append(gr.myGroups[groupPublicKey],
		&membership{
			signer:  signer,
			channel: channel,
		})
}

// UnregisterGroup removes a group from myGroup array by a public key
func (gr *GroupRegistry) UnregisterGroup(groupPublicKey string) {

	gr.mutex.Lock()
	defer gr.mutex.Unlock()

	delete(gr.myGroups, groupPublicKey)
}

func (gr *GroupRegistry) getGroup(groupPublicKey string) []*membership {
	gr.mutex.Lock()
	defer gr.mutex.Unlock()

	return gr.myGroups[groupPublicKey]
}

// RemoveExpiredGroups lookup for groups to be removed.
func (gr *GroupRegistry) RemoveExpiredGroups() {
	gr.mutex.Lock()
	defer gr.mutex.Unlock()

	for publicKey := range gr.myGroups {
		isGroupEligibleForRemoval, err := gr.relayChain.IsGroupEligibleForRemoval([]byte(publicKey))

		if err != nil {
			fmt.Fprintf(os.Stderr, "Group check eligibility failed: [%v]\n", err)
		}

		if isGroupEligibleForRemoval {
			gr.UnregisterGroup(publicKey)
		}
	}
}
