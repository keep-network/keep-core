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

// UnregisterStaleGroups lookup for groups that have been marked as stale
// on-chain. A stale group is a group that has expired and a certain time passed
// after the group expiration. This guarantees the group will not be selected to
// a new operation and it cannot have an ongoing operation for which it could be
// selected before it expired. Such a group can be safely removed from the registry
// and archived in the underlying storage.
func (gr *Groups) UnregisterStaleGroups() error {
	gr.mutex.Lock()
	defer gr.mutex.Unlock()

	for publicKey := range gr.myGroups {
		publicKeyBytes, err := groupKeyFromString(publicKey)
		if err != nil {
			return fmt.Errorf("error occured while decoding public key into bytes [%v]", err)
		}

		isStaleGroup, err := gr.relayChain.IsStaleGroup(publicKeyBytes)
		if err != nil {
			return fmt.Errorf("staling group eligibility check has failed: [%v]", err)
		}

		if isStaleGroup {
			err = gr.storage.archive(publicKey)
			if err != nil {
				return fmt.Errorf("group archiving has failed: [%v]", err)
			}

			delete(gr.myGroups, publicKey)
		}
	}

	return nil
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

		gr.myGroups[groupPublicKey] = append(gr.myGroups[groupPublicKey], membership)
	}

	for group, memberships := range gr.myGroups {
		fmt.Fprintf(os.Stdout, "Group [%v] was loaded with member IDs [", group)
		for idx, membership := range memberships {
			if (len(memberships) - 1) != idx {
				fmt.Fprintf(os.Stdout, "%v, ", membership.Signer.MemberID())
			} else {
				fmt.Fprintf(os.Stdout, "%v]\n", membership.Signer.MemberID())
			}
		}
	}

	return nil
}

func groupKeyToString(groupKey []byte) string {
	return hex.EncodeToString(groupKey)
}

func groupKeyFromString(groupKey string) ([]byte, error) {
	return hex.DecodeString(groupKey)
}
