package dkg

import (
	"fmt"
	"math/big"
	"strconv"
	"testing"

	"github.com/bnb-chain/tss-lib/ecdsa/keygen"
	"github.com/bnb-chain/tss-lib/tss"
	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/chain/local_v1"
	"github.com/keep-network/keep-core/pkg/internal/testutils"
	"github.com/keep-network/keep-core/pkg/operator"
	"github.com/keep-network/keep-core/pkg/protocol/group"
)

func TestShouldAcceptMessage(t *testing.T) {
	groupSize := 5
	honestThreshold := 3

	localChain := local_v1.Connect(groupSize, honestThreshold)

	operatorsAddresses := make([]chain.Address, groupSize)
	operatorsPublicKeys := make([][]byte, groupSize)
	for i := range operatorsAddresses {
		_, operatorPublicKey, err := operator.GenerateKeyPair(
			local_v1.DefaultCurve,
		)
		if err != nil {
			t.Fatal(err)
		}

		operatorAddress, err := localChain.Signing().PublicKeyToAddress(
			operatorPublicKey,
		)
		if err != nil {
			t.Fatal(err)
		}

		operatorsAddresses[i] = operatorAddress
		operatorsPublicKeys[i] = operator.MarshalUncompressed(operatorPublicKey)
	}

	tests := map[string]struct {
		senderID         group.MemberIndex
		senderPublicKey  []byte
		activeMembersIDs []group.MemberIndex
		expectedResult   bool
	}{
		"message from another valid and operating member": {
			senderID:         group.MemberIndex(2),
			senderPublicKey:  operatorsPublicKeys[1],
			activeMembersIDs: []group.MemberIndex{1, 2, 3, 4, 5},
			expectedResult:   true,
		},
		"message from another valid but non-operating member": {
			senderID:         group.MemberIndex(2),
			senderPublicKey:  operatorsPublicKeys[1],
			activeMembersIDs: []group.MemberIndex{1, 3, 4, 5}, // 2 is inactive
			expectedResult:   false,
		},
		"message from self": {
			senderID:         group.MemberIndex(1),
			senderPublicKey:  operatorsPublicKeys[0],
			activeMembersIDs: []group.MemberIndex{1, 2, 3, 4, 5},
			expectedResult:   false,
		},
		"message from another invalid member": {
			senderID:         group.MemberIndex(2),
			senderPublicKey:  operatorsPublicKeys[3],
			activeMembersIDs: []group.MemberIndex{1, 2, 3, 4, 5},
			expectedResult:   false,
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			membershipValdator := group.NewMembershipValidator(
				&testutils.MockLogger{},
				operatorsAddresses,
				localChain.Signing(),
			)

			member := newMember(
				&testutils.MockLogger{},
				group.MemberIndex(1),
				groupSize,
				groupSize-honestThreshold,
				membershipValdator,
				"1",
				&keygen.LocalPreParams{},
				1,
			)

			filter := member.inactiveMemberFilter()
			for _, activeMemberID := range test.activeMembersIDs {
				filter.MarkMemberAsActive(activeMemberID)
			}
			filter.FlushInactiveMembers()

			result := member.shouldAcceptMessage(test.senderID, test.senderPublicKey)

			testutils.AssertBoolsEqual(
				t,
				"result from message validator",
				test.expectedResult,
				result,
			)
		})
	}
}

func TestGenerateTssPartiesIDs(t *testing.T) {
	thisMemberID := group.MemberIndex(2)
	groupMembersIDs := []group.MemberIndex{
		group.MemberIndex(1),
		group.MemberIndex(2),
		group.MemberIndex(3),
		group.MemberIndex(4),
		group.MemberIndex(5),
	}

	thisTssPartyID, groupTssPartiesIDs := generateTssPartiesIDs(
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

	tssPartyID := resolveSortedTssPartyID(tssParameters, memberID)

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
