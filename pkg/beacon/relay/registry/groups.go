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

	groupPublicKey := groupKeyToString(signer.GroupPublicKeyBytes())

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

	return gr.myGroups[groupKeyToString(groupPublicKey)]
}

// UnregisterDeletedGroups lookup for groups to be removed.
func (gr *Groups) UnregisterDeletedGroups() {
	gr.mutex.Lock()
	defer gr.mutex.Unlock()

	for publicKey := range gr.myGroups {
		publicKeyBytes, err := groupKeyFromString(publicKey)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error occured while decoding public key into bytes [%v]\n", err)
		}

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
		groupPublicKey := groupKeyToString(membership.Signer.GroupPublicKeyBytes())
		fmt.Fprintf(os.Stdout, "Membership: [%v] was loaded to a group: [%v]\n",
			membership.Signer.MemberID(),
			groupPublicKey,
		)

		gr.myGroups[groupPublicKey] = append(gr.myGroups[groupPublicKey], membership)
	}

	return nil
}

func groupKeyToString(groupKey []byte) string {
	return hex.EncodeToString(groupKey)
}

func groupKeyFromString(groupKey string) ([]byte, error) {
	return hex.DecodeString(groupKey)
}
