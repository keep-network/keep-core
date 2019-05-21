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
}

// Membership represents a member of a group
type Membership struct {
	Signer      *dkg.ThresholdSigner
	ChannelName string
}

// NewGroups returns an empty GroupRegistry.
func NewGroups(
	relayChain relaychain.GroupRegistrationInterface,
) *Groups {
	return &Groups{
		myGroups:   make(map[string][]*Membership),
		relayChain: relayChain,
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

	groupPublicKey := string(signer.GroupPublicKeyBytes())

	membership := &Membership{
		Signer:      signer,
		ChannelName: channelName,
	}

	membershipBytes, err := membership.Marshal()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Marshalling of the membership failed: [%v]\n", err)
	}

	fileStorage := storage.NewFileStorage()
	if fileStorage != nil {
		hexGroupPublicKey := hex.EncodeToString(signer.GroupPublicKeyBytes())
		fileStorage.Save(membershipBytes, "/membership_"+hexGroupPublicKey)
	} else {
		fmt.Fprintf(os.Stderr, "An error occured while retrieving a config path")
	}

	gr.myGroups[groupPublicKey] = append(gr.myGroups[groupPublicKey], membership)
}

// GetGroup gets a group by a groupPublicKey
func (gr *Groups) GetGroup(groupPublicKey []byte) []*Membership {
	gr.mutex.Lock()
	defer gr.mutex.Unlock()

	return gr.myGroups[string(groupPublicKey)]
}

// UnregisterDeletedGroups lookup for groups to be removed.
func (gr *Groups) UnregisterDeletedGroups() {
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
