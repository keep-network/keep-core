package registry

import (
	"encoding/hex"
	"fmt"
	"sync"

	relaychain "github.com/keep-network/keep-core/pkg/beacon/relay/chain"
	"github.com/keep-network/keep-core/pkg/beacon/relay/dkg"

	"github.com/keep-network/keep-common/pkg/persistence"
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
func (g *Groups) RegisterGroup(
	signer *dkg.ThresholdSigner,
	channelName string,
) error {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	groupPublicKey := groupKeyToString(signer.GroupPublicKeyBytes())

	membership := &Membership{
		Signer:      signer,
		ChannelName: channelName,
	}

	err := g.storage.save(membership)
	if err != nil {
		return fmt.Errorf("could not persist membership to the storage: [%v]", err)
	}

	g.myGroups[groupPublicKey] = append(g.myGroups[groupPublicKey], membership)

	return nil
}

// GetGroup gets a group by a groupPublicKey
func (g *Groups) GetGroup(groupPublicKey []byte) []*Membership {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	return g.myGroups[groupKeyToString(groupPublicKey)]
}

// UnregisterStaleGroups lookup for groups that have been marked as stale
// on-chain. A stale group is a group that has expired and a certain time passed
// after the group expiration. This guarantees the group will not be selected to
// a new operation and it cannot have an ongoing operation for which it could be
// selected before it expired. Such a group can be safely removed from the registry
// and archived in the underlying storage.
func (g *Groups) UnregisterStaleGroups() {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	for publicKey := range g.myGroups {
		publicKeyBytes, err := groupKeyFromString(publicKey)
		if err != nil {
			logger.Errorf(
				"error occured while decoding public key into bytes: [%v]",
				err,
			)
		}

		isStaleGroup, err := g.relayChain.IsStaleGroup(publicKeyBytes)
		if err != nil {
			logger.Errorf("stale group check has failed: [%v]", err)
		}

		if isStaleGroup {
			err = g.storage.archive(publicKey)
			if err != nil {
				logger.Errorf("group archiving has failed: [%v]", err)
			}

			delete(g.myGroups, publicKey)
		}
	}
}

// LoadExistingGroups iterates over all stored memberships on disk and loads them
// into memory
func (g *Groups) LoadExistingGroups() {
	g.myGroups = make(map[string][]*Membership)

	membershipsChannel, errorsChannel := g.storage.readAll()

	// Two goroutines read from memberships and errors channels and either
	// adds memberships to the group registry or outputs an error to stderr.
	// The reason for using two goroutines at the same time - one for
	// memberships and one for errors is because channels do not have to be
	// buffered and we do not know in what order information is written to
	// channels.
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		for membership := range membershipsChannel {
			groupPublicKey := groupKeyToString(
				membership.Signer.GroupPublicKeyBytes(),
			)
			g.myGroups[groupPublicKey] = append(
				g.myGroups[groupPublicKey],
				membership,
			)
		}

		wg.Done()
	}()

	go func() {
		for err := range errorsChannel {
			logger.Errorf(
				"could not load membership from disk: [%v]",
				err,
			)
		}

		wg.Done()
	}()

	wg.Wait()

	g.printMemberships()
}

func (g *Groups) printMemberships() {
	for group, memberships := range g.myGroups {
		memberLog := fmt.Sprintf("group [0x%v] loaded with members: [", group)
		for idx, membership := range memberships {
			if (len(memberships) - 1) != idx {
				memberLog += fmt.Sprintf("%v, ", membership.Signer.MemberID())
			} else {
				memberLog += fmt.Sprintf("%v]", membership.Signer.MemberID())
			}
		}

		logger.Infof(memberLog)
	}
}

func groupKeyToString(groupKey []byte) string {
	return hex.EncodeToString(groupKey)
}

func groupKeyFromString(groupKey string) ([]byte, error) {
	return hex.DecodeString(groupKey)
}
