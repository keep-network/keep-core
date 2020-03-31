package group

import (
	"crypto/ecdsa"
	"encoding/hex"

	relaychain "github.com/keep-network/keep-core/pkg/beacon/relay/chain"
	"github.com/keep-network/keep-core/pkg/chain"
)

// MembershipValidator operates on a group selection result and lets to
// validate one's membership based on the provided public key.
//
// Validator is used to filter out messages from parties not selected to
// the group. It is also used to confirm the position in the group of
// a party that was selected. This is used to validate messages sent by that
// party to all other group members.
type MembershipValidator struct {
	members map[string][]int // staker address -> staker positions in group
	signing chain.Signing
}

// NewMembershipValidator creates a validator for the provided group selection
// result.
func NewMembershipValidator(
	selected []relaychain.StakerAddress,
	signing chain.Signing,
) *MembershipValidator {
	members := make(map[string][]int)
	for position, address := range selected {
		addressAsString := hex.EncodeToString(address)
		positions, ok := members[addressAsString]
		if ok {
			positions = append(positions, position)
			members[addressAsString] = positions
		} else {
			members[addressAsString] = []int{position}
		}
	}

	return &MembershipValidator{
		members: members,
		signing: signing,
	}
}

// IsInGroup returns true if party with the given public key has been selected
// to the group. Otherwise, function returns false.
func (mv *MembershipValidator) IsInGroup(publicKey *ecdsa.PublicKey) bool {
	authorAddress := hex.EncodeToString(
		mv.signing.PublicKeyToAddress(*publicKey),
	)
	_, isInGroup := mv.members[authorAddress]
	return isInGroup
}

// IsSelectedAtIndex returns true if party with the given public key has been
// selected to the group at the given position. If the position does not match
// function returns false. The same happens when the party was not selected
// to the group.
func (mv *MembershipValidator) IsSelectedAtIndex(
	index int,
	publicKey *ecdsa.PublicKey,
) bool {
	authorAddress := hex.EncodeToString(
		mv.signing.PublicKeyToAddress(*publicKey),
	)
	positions, isInGroup := mv.members[authorAddress]

	if !isInGroup {
		return false
	}

	for _, position := range positions {
		if index == position {
			return true
		}
	}

	return false
}
