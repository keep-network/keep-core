package registry

import (
	"bytes"
	"encoding/hex"
	"math/big"
	"reflect"
	"testing"

	bn256 "github.com/ethereum/go-ethereum/crypto/bn256/cloudflare"
	"github.com/keep-network/keep-common/pkg/persistence"
	"github.com/keep-network/keep-core/pkg/beacon/relay/dkg"
	"github.com/keep-network/keep-core/pkg/beacon/relay/event"
	"github.com/keep-network/keep-core/pkg/beacon/relay/group"
	chainLocal "github.com/keep-network/keep-core/pkg/chain/local"
	"github.com/keep-network/keep-core/pkg/subscription"
)

var (
	channelName1 = "test_channel1"
	channelName2 = "test_channel2"

	persistenceMock = &persistenceHandleMock{}

	signer1 = dkg.NewThresholdSigner(
		group.MemberIndex(1),
		new(bn256.G2).ScalarBaseMult(big.NewInt(10)),
		big.NewInt(1),
	)
	signer2 = dkg.NewThresholdSigner(
		group.MemberIndex(2),
		new(bn256.G2).ScalarBaseMult(big.NewInt(20)),
		big.NewInt(2),
	)
	signer3 = dkg.NewThresholdSigner(
		group.MemberIndex(3),
		new(bn256.G2).ScalarBaseMult(big.NewInt(30)),
		big.NewInt(3),
	)
	signer4 = dkg.NewThresholdSigner(
		group.MemberIndex(3),
		new(bn256.G2).ScalarBaseMult(big.NewInt(20)),
		big.NewInt(2),
	)
)

func TestRegisterGroup(t *testing.T) {
	chain := chainLocal.Connect(5, 3, big.NewInt(200)).ThresholdRelay()

	gr := NewGroupRegistry(chain, persistenceMock)

	gr.RegisterGroup(signer1, channelName1)

	actual := gr.GetGroup(signer1.GroupPublicKeyBytesCompressed())

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

func TestLoadGroup(t *testing.T) {
	chain := chainLocal.Connect(5, 3, big.NewInt(200)).ThresholdRelay()
	gr := NewGroupRegistry(chain, persistenceMock)

	if len(gr.myGroups) != 0 {
		t.Fatalf(
			"Unexpected number of group memberships at a Keep Node start \nExpected: [%+v]\nActual:   [%+v]",
			0,
			len(gr.myGroups),
		)
	}

	gr.LoadExistingGroups()

	if len(gr.myGroups) != 2 {
		t.Fatalf(
			"Unexpected number of group memberships \nExpected: [%+v]\nActual:   [%+v]",
			2,
			len(gr.myGroups),
		)
	}

	expectedMembership1 := &Membership{
		Signer:      signer1,
		ChannelName: channelName1,
	}
	actualMembership1 := gr.GetGroup(signer1.GroupPublicKeyBytesCompressed())[0]
	if !reflect.DeepEqual(expectedMembership1, actualMembership1) {
		t.Errorf("\nexpected: %v\nactual:   %v", expectedMembership1, actualMembership1)
	}

	expectedMembership2 := &Membership{
		Signer:      signer2,
		ChannelName: channelName2,
	}
	actualMembership2 := gr.GetGroup(signer2.GroupPublicKeyBytesCompressed())[0]
	if !reflect.DeepEqual(expectedMembership2, actualMembership2) {
		t.Errorf("\nexpected: %v\nactual:   %v", expectedMembership2, actualMembership2)
	}
}

func TestUnregisterStaleGroups(t *testing.T) {
	mockChain := &mockGroupRegistrationInterface{
		groupsToRemove: [][]byte{},
	}

	gr := NewGroupRegistry(mockChain, persistenceMock)

	gr.RegisterGroup(signer1, channelName1)
	gr.RegisterGroup(signer2, channelName1)
	gr.RegisterGroup(signer3, channelName1)

	mockChain.markAsStale(signer2.GroupPublicKeyBytesCompressed())

	gr.UnregisterStaleGroups()

	group1 := gr.GetGroup(signer1.GroupPublicKeyBytesCompressed())
	if group1 == nil {
		t.Fatalf("Expecting a group, but nil was returned instead")
	}

	group2 := gr.GetGroup(signer2.GroupPublicKeyBytesCompressed())
	if group2 != nil {
		t.Fatalf("Group2 was expected to be unregistered, but is still present")
	}
	if len(persistenceMock.archivedGroups) != 1 ||
		persistenceMock.archivedGroups[0] != hex.EncodeToString(signer2.GroupPublicKeyBytesCompressed()) {
		t.Fatalf("Group2 was expected to be archived")
	}

	group3 := gr.GetGroup(signer3.GroupPublicKeyBytesCompressed())
	if group3 == nil {
		t.Fatalf("Expecting a group, but nil was returned instead")
	}

}

type mockGroupRegistrationInterface struct {
	groupsToRemove [][]byte
}

func (mgri *mockGroupRegistrationInterface) markAsStale(publicKey []byte) {
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

type persistenceHandleMock struct {
	archivedGroups []string
}

func (phm *persistenceHandleMock) Save(data []byte, directory string, name string) error {
	// noop
	return nil
}

func (phm *persistenceHandleMock) ReadAll() (<-chan persistence.DataDescriptor, <-chan error) {
	membershipBytes1, _ := (&Membership{
		Signer:      signer1,
		ChannelName: channelName1,
	}).Marshal()

	membershipBytes2, _ := (&Membership{
		Signer:      signer2,
		ChannelName: channelName2,
	}).Marshal()

	membershipBytes3, _ := (&Membership{
		Signer:      signer4,
		ChannelName: channelName2,
	}).Marshal()

	outputData := make(chan persistence.DataDescriptor, 3)
	outputErrors := make(chan error)

	outputData <- &testDataDescriptor{"1", "dir", membershipBytes1}
	outputData <- &testDataDescriptor{"2", "dir", membershipBytes2}
	outputData <- &testDataDescriptor{"3", "dir", membershipBytes3}

	close(outputData)
	close(outputErrors)

	return outputData, outputErrors
}

func (phm *persistenceHandleMock) Archive(directory string) error {
	phm.archivedGroups = append(phm.archivedGroups, directory)

	return nil
}

type testDataDescriptor struct {
	name      string
	directory string
	content   []byte
}

func (tdd *testDataDescriptor) Name() string {
	return tdd.name
}

func (tdd *testDataDescriptor) Directory() string {
	return tdd.directory
}

func (tdd *testDataDescriptor) Content() ([]byte, error) {
	return tdd.content, nil
}
