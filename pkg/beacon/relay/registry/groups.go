package registry

import (
	"fmt"
	"os"
	"sync"

	"encoding/hex"

	relaychain "github.com/keep-network/keep-core/pkg/beacon/relay/chain"
	"github.com/keep-network/keep-core/pkg/beacon/relay/dkg"
	"github.com/keep-network/keep-core/pkg/storage"
)

// Groups represents a collection of Keep groups in which the given
// client is a member.
type Groups struct {
	mutex sync.Mutex

	myGroups map[string][]*Membership

	relayChain relaychain.GroupRegistrationInterface

	storage storage.Storage
}

// Membership represents a member of a group
type Membership struct {
	Signer      *dkg.ThresholdSigner
	ChannelName string
}

// NewGroupRegistry returns an empty GroupRegistry.
func NewGroupRegistry(
	relayChain relaychain.GroupRegistrationInterface,
	storage storage.Storage,
) *Groups {
	return &Groups{
		myGroups:   make(map[string][]*Membership),
		relayChain: relayChain,
		storage:    storage,
	}
}

// RegisterGroup registers that a group was successfully created by the given
// groupPublicKey.
func (gr *Groups) RegisterGroup(
	signer *dkg.ThresholdSigner,
	channelName string,
) {

	gr.mutex.Lock()
	defer gr.mutex.Unlock()

	membership := &Membership{
		Signer:      signer,
		ChannelName: channelName,
	}

	membershipBytes, err := membership.Marshal()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Marshalling of the membership failed: [%v]\n", err)
		return
	}
	groupPublicKey := hex.EncodeToString(signer.GroupPublicKeyBytes())
	gr.storage.Save(membershipBytes, "/membership_"+groupPublicKey)

	gr.myGroups[groupPublicKey] = append(gr.myGroups[groupPublicKey], membership)
}

// GetGroup gets a group by a groupPublicKey
func (gr *Groups) GetGroup(groupPublicKey []byte) []*Membership {
	gr.mutex.Lock()
	defer gr.mutex.Unlock()

	return gr.myGroups[hex.EncodeToString(groupPublicKey)]
}

// UnregisterDeletedGroups lookup for groups to be removed.
func (gr *Groups) UnregisterDeletedGroups() {
	gr.mutex.Lock()
	defer gr.mutex.Unlock()

	for publicKey := range gr.myGroups {
		publicKeyBytes, _ := hex.DecodeString(publicKey)
		isStaleGroup, err := gr.relayChain.IsStaleGroup(publicKeyBytes)

		if err != nil {
			fmt.Fprintf(os.Stderr, "Group removal eligibility check failed: [%v]\n", err)
		}

		if isStaleGroup {
			delete(gr.myGroups, publicKey)
		}
	}
}

// LoadExistingGroups iterates over all stored memberships on disk and loads them
// into memory
func (gr *Groups) LoadExistingGroups() error {
	storedMemberships := gr.storage.ReadAll()

	for _, storedMembership := range storedMemberships {
		membership := &Membership{}
		err := membership.Unmarshal(storedMembership)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error occured while unmarshalling a membership: [%v]\n", err)
			gr.myGroups = make(map[string][]*Membership)
			return err
		}

		groupPublicKey := hex.EncodeToString(membership.Signer.GroupPublicKeyBytes())
		gr.myGroups[groupPublicKey] = append(gr.myGroups[groupPublicKey], membership)
	}

	return nil
}
