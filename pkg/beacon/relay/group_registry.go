package relay

import (
	"fmt"
	"os"
	"sync"

	relaychain "github.com/keep-network/keep-core/pkg/beacon/relay/chain"
	"github.com/keep-network/keep-core/pkg/beacon/relay/dkg"
	"github.com/keep-network/keep-core/pkg/net"
)

// GroupRegistry represents a collection of Keep groups in which the given
// client is a member.
type GroupRegistry struct {
	mutex sync.Mutex

	myGroups map[string][]*Membership

	relayChain relaychain.GroupRegistrationInterface
}

// Membership represents a member of a group
type Membership struct {
	signer  *dkg.ThresholdSigner
	channel net.BroadcastChannel
}

// NewGroupRegistry returns an empty GroupRegistry.
func NewGroupRegistry(
	relayChain relaychain.GroupRegistrationInterface,
) *GroupRegistry {
	return &GroupRegistry{
		myGroups:   make(map[string][]*Membership),
		relayChain: relayChain,
	}
}

// RegisterGroup registers that a group was successfully created by the given
// groupPublicKey.
func (gr *GroupRegistry) RegisterGroup(
	signer *dkg.ThresholdSigner,
	channel net.BroadcastChannel,
) {

	gr.mutex.Lock()
	defer gr.mutex.Unlock()

	groupPublicKey := string(signer.GroupPublicKeyBytes())

	gr.myGroups[groupPublicKey] = append(gr.myGroups[groupPublicKey],
		&Membership{
			signer:  signer,
			channel: channel,
		})
}

// GetGroup gets a group by a groupPublicKey
func (gr *GroupRegistry) GetGroup(groupPublicKey []byte) []*Membership {
	gr.mutex.Lock()
	defer gr.mutex.Unlock()

	return gr.myGroups[string(groupPublicKey)]
}

// UnregisterDeletedGroups lookup for groups to be removed.
// Group is removed if it is considered as stale on-chain.
func (gr *GroupRegistry) UnregisterDeletedGroups() {
	gr.mutex.Lock()
	defer gr.mutex.Unlock()

	for publicKey := range gr.myGroups {
		publicKeyBytes := []byte(publicKey)
		isStaleGroup, err := gr.relayChain.IsStaleGroup(publicKeyBytes)

		if err != nil {
			fmt.Fprintf(os.Stderr, "Group removal eligibility check failed: [%v]\n", err)
		}

		if isStaleGroup {
			delete(gr.myGroups, publicKey)
			fmt.Printf("Unregistering a stale group [%+v]\n", publicKeyBytes)
		}
	}
}
