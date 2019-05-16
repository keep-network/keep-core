package registry

import (
	"bytes"
	"math/big"
	"sync"
	"testing"

	bn256 "github.com/ethereum/go-ethereum/crypto/bn256/cloudflare"
	"github.com/keep-network/keep-core/pkg/beacon/relay/dkg"
	"github.com/keep-network/keep-core/pkg/beacon/relay/event"
	"github.com/keep-network/keep-core/pkg/beacon/relay/group"
	chainLocal "github.com/keep-network/keep-core/pkg/chain/local"
	netLocal "github.com/keep-network/keep-core/pkg/net/local"
	"github.com/keep-network/keep-core/pkg/subscription"
)

func TestRegisterGroup(t *testing.T) {
	signer := dkg.NewThresholdSigner(
		group.MemberIndex(2),
		new(bn256.G2).ScalarBaseMult(big.NewInt(10)),
		big.NewInt(1),
	)

	gr := &Groups{
		mutex:      sync.Mutex{},
		myGroups:   make(map[string][]*Membership),
		relayChain: chainLocal.Connect(5, 3, big.NewInt(200)).ThresholdRelay(),
	}

	networkProvider := netLocal.Connect()
	channel, err := networkProvider.ChannelFor("testChannel")
	if err != nil {
		t.Fatal(err)
	}

	gr.RegisterGroup(signer, channel)

	actual := gr.GetGroup(signer.GroupPublicKeyBytes())

	if actual == nil {
		t.Fatalf(
			"Expecting a group, but nil was returned instead",
		)
	}

	if len(actual) != 1 {
		t.Fatalf(
			"Unexpected number of group memberships \nExpected: [%+v]\nActual:   [%+v]",
			1,
			len(actual),
		)
	}
}

func TestUnregisterStaleGroups(t *testing.T) {
	mockChain := &mockGroupRegistrationInterface{
		groupsToRemove: [][]byte{},
	}

	gr := &Groups{
		mutex:      sync.Mutex{},
		myGroups:   make(map[string][]*Membership),
		relayChain: mockChain,
	}

	networkProvider := netLocal.Connect()
	channel, err := networkProvider.ChannelFor("test")
	if err != nil {
		t.Fatal(err)
	}

	signer1 := dkg.NewThresholdSigner(
		group.MemberIndex(1),
		new(bn256.G2).ScalarBaseMult(big.NewInt(10)),
		big.NewInt(1),
	)
	signer2 := dkg.NewThresholdSigner(
		group.MemberIndex(2),
		new(bn256.G2).ScalarBaseMult(big.NewInt(20)),
		big.NewInt(2),
	)
	signer3 := dkg.NewThresholdSigner(
		group.MemberIndex(3),
		new(bn256.G2).ScalarBaseMult(big.NewInt(30)),
		big.NewInt(3),
	)

	gr.RegisterGroup(signer1, channel)
	gr.RegisterGroup(signer2, channel)
	gr.RegisterGroup(signer3, channel)

	mockChain.markForRemoval(signer2.GroupPublicKeyBytes())

	gr.UnregisterDeletedGroups()

	group1 := gr.GetGroup(signer1.GroupPublicKeyBytes())
	if group1 == nil {
		t.Fatalf(
			"Expecting a group, but nil was returned instead",
		)
	}

	group2 := gr.GetGroup(signer2.GroupPublicKeyBytes())
	if group2 != nil {
		t.Fatalf(
			"Group2 was expected to be unregistered, but is still present",
		)
	}

	group3 := gr.GetGroup(signer3.GroupPublicKeyBytes())
	if group3 == nil {
		t.Fatalf(
			"Expecting a group, but nil was returned instead",
		)
	}

}

type mockGroupRegistrationInterface struct {
	groupsToRemove [][]byte
}

func (mgri *mockGroupRegistrationInterface) markForRemoval(publicKey []byte) {
	mgri.groupsToRemove = append(mgri.groupsToRemove, publicKey)
}

func (mgri *mockGroupRegistrationInterface) OnGroupRegistered(
	func(groupRegistration *event.GroupRegistration),
) (subscription.EventSubscription, error) {
	panic("not implemented")
}

func (mgri *mockGroupRegistrationInterface) IsStaleGroup(groupPublicKey []byte) (bool, error) {
	for _, groupToRemove := range mgri.groupsToRemove {
		if bytes.Compare(groupToRemove, groupPublicKey) == 0 {
			return true, nil
		}
	}
	return false, nil
}
