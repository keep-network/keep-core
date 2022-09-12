package signing

import (
	"fmt"
	"github.com/bnb-chain/tss-lib/tss"
	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/chain/local_v1"
	"github.com/keep-network/keep-core/pkg/internal/tecdsatest"
	"github.com/keep-network/keep-core/pkg/internal/testutils"
	"github.com/keep-network/keep-core/pkg/operator"
	"github.com/keep-network/keep-core/pkg/protocol/group"
	"github.com/keep-network/keep-core/pkg/tecdsa"
	"math/big"
	"testing"
)

func TestShouldAcceptMessage(t *testing.T) {
	groupSize := 5
	honestThreshold := 3

	localChain := local_v1.Connect(groupSize, honestThreshold)

	testData, err := tecdsatest.LoadPrivateKeyShareTestFixtures(1)
	if err != nil {
		t.Fatalf("failed to load test data: [%v]", err)
	}

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
				big.NewInt(100),
				tecdsa.NewPrivateKeyShare(testData[0]),
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

func TestIdentityConverter_MemberIndexToTssPartyID(t *testing.T) {
	converter := &identityConverter{keys: []*big.Int{
		big.NewInt(301),
		big.NewInt(303),
		big.NewInt(304),
		big.NewInt(305),
	}}
	memberID := group.MemberIndex(2)

	tssPartyID := converter.MemberIndexToTssPartyID(memberID)

	testutils.AssertStringsEqual(
		t,
		"ID of the TSS party ID",
		tssPartyID.Id,
		"303",
	)

	testutils.AssertBytesEqual(
		t,
		tssPartyID.Key,
		big.NewInt(303).Bytes(),
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

func TestIdentityConverter_MemberIndexToTssPartyIDKey(t *testing.T) {
	converter := &identityConverter{keys: []*big.Int{
		big.NewInt(301),
		big.NewInt(303),
		big.NewInt(304),
		big.NewInt(305),
	}}
	memberID := group.MemberIndex(2)

	key := converter.MemberIndexToTssPartyIDKey(memberID)

	testutils.AssertBigIntsEqual(
		t,
		"key of the TSS party ID",
		big.NewInt(303),
		key,
	)
}

func TestIdentityConverter_TssPartyIDToMemberIndex(t *testing.T) {
	converter := &identityConverter{keys: []*big.Int{
		big.NewInt(301),
		big.NewInt(303),
		big.NewInt(304),
		big.NewInt(305),
	}}
	partyID := tss.NewPartyID("303", "member-2", big.NewInt(303))

	memberID := converter.TssPartyIDToMemberIndex(partyID)

	testutils.AssertIntsEqual(t, "member ID", 2, int(memberID))
}