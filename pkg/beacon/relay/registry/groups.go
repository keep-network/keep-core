package registry

import (
	"encoding/hex"
	"fmt"
	"os"
	"sync"

	relaychain "github.com/keep-network/keep-core/pkg/beacon/relay/chain"
	"github.com/keep-network/keep-core/pkg/beacon/relay/dkg"

	"github.com/keep-network/keep-core/pkg/persistence"
)

// Groups represents a collection of Keep groups in which the given
// client is a member.
type Groups struct {
	mutex sync.Mutex

	myGroups map[string][]*Membership

	relayChain relaychain.GroupRegistrationInterface

	storage storage
}

// Membership represents a member of a group
type Membership struct {
	Signer      *dkg.ThresholdSigner
	ChannelName string
}

// NewGroupRegistry returns an empty GroupRegistry.
func NewGroupRegistry(
	relayChain relaychain.GroupRegistrationInterface,
	persistence persistence.Handle,
) *Groups {
	return &Groups{
		myGroups:   make(map[string][]*Membership),
		relayChain: relayChain,
		storage:    newStorage(persistence),
		mutex:      sync.Mutex{},
	}
}

// RegisterGroup registers that a group was successfully created by the given
// groupPublicKey.
func (gr *Groups) RegisterGroup(
	signer *dkg.ThresholdSigner,
	channelName string,
) error {
	gr.mutex.Lock()
	defer gr.mutex.Unlock()

	groupPublicKey := hex.EncodeToString(signer.GroupPublicKeyBytes())

	membership := &Membership{
		Signer:      signer,
		ChannelName: channelName,
	}

	err := gr.storage.save(membership)
	if err != nil {
		return fmt.Errorf("could not persist membership to the storage: [%v]", err)
	}

	gr.myGroups[groupPublicKey] = append(gr.myGroups[groupPublicKey], membership)

	return nil
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
	memberships, err := gr.storage.readAll()
	if err != nil {
		gr.myGroups = make(map[string][]*Membership)
		return err
	}

	for _, membership := range memberships {
		groupPublicKey := hex.EncodeToString(membership.Signer.GroupPublicKeyBytes())
		fmt.Printf("groupPublicKey: [%s]", groupPublicKey)
		gr.myGroups[groupPublicKey] = append(gr.myGroups[groupPublicKey], membership)
	}

	return nil
}
