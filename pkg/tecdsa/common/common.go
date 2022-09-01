// Package common holds some common tools that can be used across multiple
// tECDSA protocols, e.g. DKG and signing.
package common

import (
	"fmt"
	"github.com/bnb-chain/tss-lib/tss"
	"github.com/keep-network/keep-core/pkg/crypto/ephemeral"
	"github.com/keep-network/keep-core/pkg/protocol/group"
	"math/big"
	"strconv"
)

// GenerateTssPartiesIDs converts group member ID to parties ID suitable for
// the TSS protocol execution.
func GenerateTssPartiesIDs(
	memberID group.MemberIndex,
	groupMembersIDs []group.MemberIndex,
) (*tss.PartyID, []*tss.PartyID) {
	var partyID *tss.PartyID
	groupPartiesIDs := make([]*tss.PartyID, len(groupMembersIDs))

	for i, groupMemberID := range groupMembersIDs {
		newPartyID := NewTssPartyIDFromMemberID(groupMemberID)

		if memberID == groupMemberID {
			partyID = newPartyID
		}

		groupPartiesIDs[i] = newPartyID
	}

	return partyID, groupPartiesIDs
}

// NewTssPartyIDFromMemberID creates a new instance of a TSS party ID using
// the given member ID. Such a created party ID has an unset index since it
// does not yet belong to a sorted parties IDs set.
func NewTssPartyIDFromMemberID(memberID group.MemberIndex) *tss.PartyID {
	return tss.NewPartyID(
		strconv.Itoa(int(memberID)),
		fmt.Sprintf("member-%v", memberID),
		MemberIDToTssPartyIDKey(memberID),
	)
}

// MemberIDToTssPartyIDKey converts a single group member ID to a key that
// can be used to create a TSS party ID.
func MemberIDToTssPartyIDKey(memberID group.MemberIndex) *big.Int {
	return big.NewInt(int64(memberID))
}

// TssPartyIDToMemberID converts a single TSS party ID to a group member ID.
func TssPartyIDToMemberID(partyID *tss.PartyID) group.MemberIndex {
	return group.MemberIndex(partyID.KeyInt().Int64())
}

// ResolveSortedTssPartyID resolves the TSS party ID for the given member ID
// based on the sorted parties IDs stored in the given TSS parameters set. Such
// a resolved party ID has an index which indicates its position in the parties
// IDs set.
func ResolveSortedTssPartyID(
	tssParameters *tss.Parameters,
	memberID group.MemberIndex,
) *tss.PartyID {
	sortedPartiesIDs := tssParameters.Parties().IDs()
	partyIDKey := MemberIDToTssPartyIDKey(memberID)
	return sortedPartiesIDs.FindByKey(partyIDKey)
}

// AggregateTssMessages takes a list of TSS messages and build an aggregate
// consisting of unencrypted broadcast part and encrypted P2P parts intended
// for specific receivers. The encryption of a specific P2P part is done using
// a symmetric key taken from the provided symmetricKeys map using the
// receiver member index as key.
//
// This function has also the following requirements regarding the input
// tssMessages list and symmetricKeys map:
// - tssMessages MUST hold 0 or 1 broadcast message
// - tssMessages MUST hold 0 or N P2P messages. If P2P messages count is >0 then:
//     - each P2P message MUST target exactly 1 unique receiver
//     - each P2P message receiver MUST have a corresponding entry in
//       the symmetricKeys map
func AggregateTssMessages(
	tssMessages []tss.Message,
	symmetricKeys map[group.MemberIndex]ephemeral.SymmetricKey,
) (
	[]byte,
	map[group.MemberIndex][]byte,
	error,
) {
	var broadcastPayload []byte
	peersPayload := make(map[group.MemberIndex][]byte)

	for _, tssMessage := range tssMessages {
		tssMessageBytes, tssMessageRouting, err := tssMessage.WireBytes()
		if err != nil {
			return nil, nil, fmt.Errorf(
				"failed to unpack TSS message: [%v]",
				err,
			)
		}

		if tssMessageRouting.IsBroadcast {
			// We expect only one broadcast message. Any other case is an error.
			if len(broadcastPayload) > 0 {
				return nil, nil, fmt.Errorf(
					"multiple TSS broadcast messages detected",
				)
			}

			broadcastPayload = tssMessageBytes
		} else {
			// We expect that each P2P message targets only a single member.
			// Any other case is an error.
			if len(tssMessageRouting.To) != 1 {
				return nil, nil, fmt.Errorf(
					"multi-receiver TSS P2P message detected",
				)
			}
			// Get the single receiver ID.
			receiverID := TssPartyIDToMemberID(tssMessageRouting.To[0])
			// Get the symmetric key with the receiver. If the symmetric key
			// cannot be found, something awful happened.
			symmetricKey, ok := symmetricKeys[receiverID]
			if !ok {
				return nil, nil, fmt.Errorf(
					"cannot get symmetric key with member [%v]",
					receiverID,
				)
			}
			// Encrypt the payload using the receiver symmetric key.
			encryptedTssMessageBytes, err := symmetricKey.Encrypt(
				tssMessageBytes,
			)
			if err != nil {
				return nil, nil, fmt.Errorf(
					"cannot encrypt TSS P2P message for member [%v]: [%v]",
					receiverID,
					err,
				)
			}

			if _, exists := peersPayload[receiverID]; exists {
				return nil, nil, fmt.Errorf(
					"duplicate TSS P2P message for member [%v]",
					receiverID,
				)
			}

			peersPayload[receiverID] = encryptedTssMessageBytes
		}
	}

	return broadcastPayload, peersPayload, nil
}