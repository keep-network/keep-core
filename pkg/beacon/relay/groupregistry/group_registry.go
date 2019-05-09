package groupregistry

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
	Signer  *dkg.ThresholdSigner
	Channel net.BroadcastChannel
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
			Signer:  signer,
			Channel: channel,
		})
}

// GetGroup gets a group by a groupPublicKey
func (gr *GroupRegistry) GetGroup(groupPublicKey []byte) []*Membership {
	gr.mutex.Lock()
	defer gr.mutex.Unlock()

	return gr.myGroups[string(groupPublicKey)]
}

// UnregisterDeletedGroups lookup for groups to be removed.
func (gr *GroupRegistry) UnregisterDeletedGroups() {
	gr.mutex.Lock()
	defer gr.mutex.Unlock()

	for publicKey := range gr.myGroups {
		publicKeyBytes := []byte(publicKey)
		isGroupRegistered, err := gr.relayChain.IsGroupRegistered(publicKeyBytes)

		if err != nil {
			fmt.Fprintf(os.Stderr, "Group removal eligibility check failed: [%v]\n", err)
		}

		if !isGroupRegistered {
			delete(gr.myGroups, publicKey)
			fmt.Printf("Unregistering a group which was removed on chain [%+v]\n", publicKeyBytes)
		}
	}
}
