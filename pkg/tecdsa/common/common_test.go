package common

import (
	"bytes"
	"fmt"
	"github.com/bnb-chain/tss-lib/tss"
	"github.com/keep-network/keep-core/pkg/crypto/ephemeral"
	"github.com/keep-network/keep-core/pkg/internal/testutils"
	"github.com/keep-network/keep-core/pkg/protocol/group"
	"math/big"
	"reflect"
	"strconv"
	"testing"
)

func TestGenerateTssPartiesIDs(t *testing.T) {
	thisMemberID := group.MemberIndex(2)
	groupMembersIDs := []group.MemberIndex{
		group.MemberIndex(1),
		group.MemberIndex(2),
		group.MemberIndex(3),
		group.MemberIndex(4),
		group.MemberIndex(5),
	}

	thisTssPartyID, groupTssPartiesIDs := GenerateTssPartiesIDs(
		thisMemberID,
		groupMembersIDs,
	)

	// Just check that the `thisTssPartyID` points to `thisMemberID`. Extensive
	// check will be done in the loop below.
	testutils.AssertBytesEqual(
		t,
		thisTssPartyID.Key,
		big.NewInt(int64(thisMemberID)).Bytes(),
	)

	testutils.AssertIntsEqual(
		t,
		"length of resulting group TSS parties IDs",
		len(groupMembersIDs),
		len(groupTssPartiesIDs),
	)

	for i, tssPartyID := range groupTssPartiesIDs {
		testutils.AssertStringsEqual(
			t,
			fmt.Sprintf("ID of the TSS party ID [%v]", i),
			tssPartyID.Id,
			strconv.Itoa(int(groupMembersIDs[i])),
		)

		testutils.AssertBytesEqual(
			t,
			tssPartyID.Key,
			big.NewInt(int64(groupMembersIDs[i])).Bytes(),
		)

		testutils.AssertStringsEqual(
			t,
			fmt.Sprintf("moniker of the TSS party ID [%v]", i),
			tssPartyID.Moniker,
			fmt.Sprintf("member-%v", groupMembersIDs[i]),
		)

		testutils.AssertIntsEqual(
			t,
			fmt.Sprintf("index of the TSS party ID [%v]", i),
			-1,
			tssPartyID.Index,
		)
	}
}

func TestNewTssPartyIDFromMemberID(t *testing.T) {
	memberID := group.MemberIndex(2)

	tssPartyID := newTssPartyIDFromMemberID(memberID)

	testutils.AssertStringsEqual(
		t,
		"ID of the TSS party ID",
		tssPartyID.Id,
		strconv.Itoa(int(memberID)),
	)

	testutils.AssertBytesEqual(
		t,
		tssPartyID.Key,
		big.NewInt(int64(memberID)).Bytes(),
	)

	testutils.AssertStringsEqual(
		t,
		"moniker of the TSS party ID",
		tssPartyID.Moniker,
		fmt.Sprintf("member-%v", memberID),
	)

	testutils.AssertIntsEqual(
		t,
		"index of the TSS party ID",
		-1,
		tssPartyID.Index,
	)
}

func TestMemberIDToTssPartyIDKey(t *testing.T) {
	memberID := group.MemberIndex(2)

	key := memberIDToTssPartyIDKey(memberID)

	testutils.AssertBigIntsEqual(
		t,
		"key of the TSS party ID",
		big.NewInt(int64(memberID)),
		key,
	)
}

func TestTssPartyIDToMemberID(t *testing.T) {
	partyID := tss.NewPartyID("2", "member-2", big.NewInt(2))

	memberID := tssPartyIDToMemberID(partyID)

	testutils.AssertIntsEqual(t, "member ID", 2, int(memberID))
}

func TestResolveSortedTssPartyID(t *testing.T) {
	groupTssPartiesIDs := []*tss.PartyID{
		tss.NewPartyID("1", "member-1", big.NewInt(1)),
		tss.NewPartyID("2", "member-2", big.NewInt(2)),
		tss.NewPartyID("3", "member-3", big.NewInt(3)),
		tss.NewPartyID("4", "member-4", big.NewInt(4)),
		tss.NewPartyID("5", "member-5", big.NewInt(5)),
	}

	tssParameters := tss.NewParameters(
		tss.EC(),
		tss.NewPeerContext(tss.SortPartyIDs(groupTssPartiesIDs)),
		groupTssPartiesIDs[0],
		len(groupTssPartiesIDs),
		2,
	)

	memberID := group.MemberIndex(2)

	tssPartyID := ResolveSortedTssPartyID(tssParameters, memberID)

	testutils.AssertStringsEqual(
		t,
		"ID of the TSS party ID",
		tssPartyID.Id,
		strconv.Itoa(int(memberID)),
	)

	testutils.AssertBytesEqual(
		t,
		tssPartyID.Key,
		big.NewInt(int64(memberID)).Bytes(),
	)

	testutils.AssertStringsEqual(
		t,
		"moniker of the TSS party ID",
		tssPartyID.Moniker,
		fmt.Sprintf("member-%v", memberID),
	)

	testutils.AssertIntsEqual(
		t,
		"index of the TSS party ID",
		int(memberID-1),
		tssPartyID.Index,
	)
}

func TestAggregateTssMessages(t *testing.T) {
	var tests = map[string]struct{
		tssMessages              []tss.Message
		symmetricKeys            map[group.MemberIndex]ephemeral.SymmetricKey
		expectedBroadcastPayload []byte
		expectedPeersPayload     map[group.MemberIndex][]byte
		expectedErr              error
	}{
		"happy path": {
			tssMessages: []tss.Message {
				newMockTssMessage([]byte{0x02}, group.MemberIndex(2)),
				newMockTssMessage([]byte{0xAA}), // broadcast message
				newMockTssMessage([]byte{0x01}, group.MemberIndex(1)),
				newMockTssMessage([]byte{0x03}, group.MemberIndex(3)),
			},
			symmetricKeys: map[group.MemberIndex]ephemeral.SymmetricKey{
				1: &mockSymmetricKey{[]byte{0x01}},
				2: &mockSymmetricKey{[]byte{0x02}},
				3: &mockSymmetricKey{[]byte{0x03}},
			},
			expectedBroadcastPayload: []byte{0xAA},
			expectedPeersPayload: map[group.MemberIndex][]byte{
				1: {0x01, 0xFF, 0x01},
				2: {0x02, 0xFF, 0x02},
				3: {0x03, 0xFF, 0x03},
			},
		},
		"only one broadcast message": {
			tssMessages: []tss.Message {
				newMockTssMessage([]byte{0xAA}), // broadcast message
			},
			expectedBroadcastPayload: []byte{0xAA},
			expectedPeersPayload: make(map[group.MemberIndex][]byte),
		},
		"only P2P messages": {
			tssMessages: []tss.Message {
				newMockTssMessage([]byte{0x02}, group.MemberIndex(2)),
				newMockTssMessage([]byte{0x01}, group.MemberIndex(1)),
				newMockTssMessage([]byte{0x03}, group.MemberIndex(3)),
			},
			symmetricKeys: map[group.MemberIndex]ephemeral.SymmetricKey{
				1: &mockSymmetricKey{[]byte{0x01}},
				2: &mockSymmetricKey{[]byte{0x02}},
				3: &mockSymmetricKey{[]byte{0x03}},
			},
			expectedBroadcastPayload: nil,
			expectedPeersPayload: map[group.MemberIndex][]byte{
				1: {0x01, 0xFF, 0x01},
				2: {0x02, 0xFF, 0x02},
				3: {0x03, 0xFF, 0x03},
			},
		},
		"multiple broadcast messages": {
			tssMessages: []tss.Message {
				newMockTssMessage([]byte{0x02}, group.MemberIndex(2)),
				newMockTssMessage([]byte{0xAA}), // broadcast message
				newMockTssMessage([]byte{0xBB}), // another broadcast message
				newMockTssMessage([]byte{0x01}, group.MemberIndex(1)),
				newMockTssMessage([]byte{0x03}, group.MemberIndex(3)),
			},
			symmetricKeys: map[group.MemberIndex]ephemeral.SymmetricKey{
				1: &mockSymmetricKey{[]byte{0x01}},
				2: &mockSymmetricKey{[]byte{0x02}},
				3: &mockSymmetricKey{[]byte{0x03}},
			},
			expectedErr: fmt.Errorf("multiple TSS broadcast messages detected"),
		},
		"P2P message with multiple receivers": {
			tssMessages: []tss.Message {
				newMockTssMessage([]byte{0x02}, group.MemberIndex(2)),
				newMockTssMessage([]byte{0xAA}), // broadcast message
				newMockTssMessage([]byte{0x01}, group.MemberIndex(1), group.MemberIndex(4)), // multiple receivers
				newMockTssMessage([]byte{0x03}, group.MemberIndex(3)),
			},
			symmetricKeys: map[group.MemberIndex]ephemeral.SymmetricKey{
				1: &mockSymmetricKey{[]byte{0x01}},
				2: &mockSymmetricKey{[]byte{0x02}},
				3: &mockSymmetricKey{[]byte{0x03}},
			},
			expectedErr: fmt.Errorf("multi-receiver TSS P2P message detected"),
		},
		"multiple P2P messages for same receiver": {
			tssMessages: []tss.Message {
				newMockTssMessage([]byte{0x02}, group.MemberIndex(2)),
				newMockTssMessage([]byte{0xAA}), // broadcast message
				newMockTssMessage([]byte{0x01}, group.MemberIndex(2)), // duplicated P2P message
				newMockTssMessage([]byte{0x03}, group.MemberIndex(3)),
			},
			symmetricKeys: map[group.MemberIndex]ephemeral.SymmetricKey{
				1: &mockSymmetricKey{[]byte{0x01}},
				2: &mockSymmetricKey{[]byte{0x02}},
				3: &mockSymmetricKey{[]byte{0x03}},
			},
			expectedErr: fmt.Errorf("duplicate TSS P2P message for member [2]"),
		},
		"missing symmetric key for receiver": {
			tssMessages: []tss.Message {
				newMockTssMessage([]byte{0x02}, group.MemberIndex(2)),
				newMockTssMessage([]byte{0xAA}), // broadcast message
				newMockTssMessage([]byte{0x01}, group.MemberIndex(1)),
				newMockTssMessage([]byte{0x03}, group.MemberIndex(3)),
			},
			symmetricKeys: map[group.MemberIndex]ephemeral.SymmetricKey{
				// missing key for member 1
				2: &mockSymmetricKey{[]byte{0x02}},
				3: &mockSymmetricKey{[]byte{0x03}},
			},
			expectedErr: fmt.Errorf("cannot get symmetric key with member [1]"),
		},
		"encryption error for receiver": {
			tssMessages: []tss.Message {
				newMockTssMessage([]byte{0x02}, group.MemberIndex(2)),
				newMockTssMessage([]byte{0xAA}), // broadcast message
				newMockTssMessage([]byte{0x01}, group.MemberIndex(1)),
				newMockTssMessage([]byte{0x03}, group.MemberIndex(3)),
			},
			symmetricKeys: map[group.MemberIndex]ephemeral.SymmetricKey{
				1: &mockSymmetricKey{[]byte{0x01}},
				2: &mockSymmetricKey{[]byte{0x02}},
				3: &mockSymmetricKey{}, // wrong key
			},
			expectedErr: fmt.Errorf("cannot encrypt TSS P2P message for member [3]: [wrong key]"),
		},
		"empty tss messages slice": {
			tssMessages: make([]tss.Message, 0),
			symmetricKeys: map[group.MemberIndex]ephemeral.SymmetricKey{
				1: &mockSymmetricKey{[]byte{0x01}},
				2: &mockSymmetricKey{[]byte{0x02}},
				3: &mockSymmetricKey{[]byte{0x03}},
			},
			expectedBroadcastPayload: nil,
			expectedPeersPayload: make(map[group.MemberIndex][]byte),
		},
		"nil tss messages slice": {
			tssMessages: nil,
			symmetricKeys: map[group.MemberIndex]ephemeral.SymmetricKey{
				1: &mockSymmetricKey{[]byte{0x01}},
				2: &mockSymmetricKey{[]byte{0x02}},
				3: &mockSymmetricKey{[]byte{0x03}},
			},
			expectedBroadcastPayload: nil,
			expectedPeersPayload: make(map[group.MemberIndex][]byte),
		},
		"nil symmetric keys map": {
			tssMessages: []tss.Message {
				newMockTssMessage([]byte{0x02}, group.MemberIndex(2)),
				newMockTssMessage([]byte{0xAA}), // broadcast message
				newMockTssMessage([]byte{0x01}, group.MemberIndex(1)),
				newMockTssMessage([]byte{0x03}, group.MemberIndex(3)),
			},
			symmetricKeys: nil,
			expectedErr: fmt.Errorf("cannot get symmetric key with member [2]"),
		},
	}
	
	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			broadcastPayload, peersPayload, err := AggregateTssMessages(
				test.tssMessages, 
				test.symmetricKeys,
			)
			
			if !bytes.Equal(test.expectedBroadcastPayload, broadcastPayload) {
				t.Errorf(
					"unexpected broadcast payload\n" +
						"expected: [%v]\n" +
						"actual:   [%v]", 
					test.expectedBroadcastPayload, 
					broadcastPayload,
				)
			}

			if !reflect.DeepEqual(test.expectedPeersPayload, peersPayload) {
				t.Errorf(
					"unexpected peers payload\n" +
						"expected: [%v]\n" +
						"actual:   [%v]",
					test.expectedPeersPayload,
					peersPayload,
				)
			}

			if !reflect.DeepEqual(test.expectedErr, err) {
				t.Errorf(
					"unexpected error\n" +
						"expected: [%v]\n" +
						"actual:   [%v]",
					test.expectedErr,
					err,
				)
			}
		})
	}
}

type mockTssMessage struct {
	bytes       []byte
	to          []*tss.PartyID
	isBroadcast bool
}

func newMockTssMessage(
	bytes []byte,
	receivers... group.MemberIndex,
) *mockTssMessage {
	var to []*tss.PartyID

	for _, receiver := range receivers {
		to = append(to, newTssPartyIDFromMemberID(receiver))
	}

	return &mockTssMessage{
		bytes:       bytes,
		to:          to,
		isBroadcast: len(to) == 0,
	}
}

func (mtm *mockTssMessage) Type() string {
	panic("not implemented")
}

func (mtm *mockTssMessage) GetTo() []*tss.PartyID {
	panic("not implemented")
}

func (mtm *mockTssMessage) GetFrom() *tss.PartyID {
	panic("not implemented")
}

func (mtm *mockTssMessage) IsBroadcast() bool {
	panic("not implemented")
}

func (mtm *mockTssMessage) IsToOldCommittee() bool {
	panic("not implemented")
}

func (mtm *mockTssMessage) IsToOldAndNewCommittees() bool {
	panic("not implemented")
}

func (mtm *mockTssMessage) WireBytes() ([]byte, *tss.MessageRouting, error) {
	return mtm.bytes, &tss.MessageRouting{
		To:          mtm.to,
		IsBroadcast: mtm.isBroadcast,
	}, nil
}

func (mtm *mockTssMessage) WireMsg() *tss.MessageWrapper {
	panic("not implemented")
}

func (mtm *mockTssMessage) String() string {
	panic("not implemented")
}

type mockSymmetricKey struct {
	key []byte
}

// Encrypt makes a fake encryption by appending 0xFF + key at the end of
// the encrypted bytes.
func (msk *mockSymmetricKey) Encrypt(bytes []byte) ([]byte, error) {
	if len(msk.key) == 0 {
		return nil, fmt.Errorf("wrong key")
	}

	return append(bytes, append([]byte{0xFF}, msk.key...)...), nil
}

func (msk *mockSymmetricKey) Decrypt(bytes []byte) ([]byte, error) {
	panic("not implemented")
}


