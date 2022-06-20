package group

import (
	"encoding/hex"
	"github.com/keep-network/keep-core/pkg/operator"

	relaychain "github.com/keep-network/keep-core/pkg/beacon/relay/chain"
	"github.com/keep-network/keep-core/pkg/chain"
)

// MembershipValidator lets to validate one's membership based on the
// provided public key.
type MembershipValidator interface {
	IsInGroup(publicKey *operator.PublicKey) bool
	IsValidMembership(memberID MemberIndex, publicKey []byte) bool
}

// StakersMembershipValidator operates on a group selection result and lets to
// validate one's membership based on the provided public key.
//
// Validator is used to filter out messages from parties not selected to
// the group. It is also used to confirm the position in the group of
// a party that was selected. This is used to validate messages sent by that
// party to all other group members.
type StakersMembershipValidator struct {
	members map[string][]int // staker address -> staker positions in group
	signing chain.Signing
}

// NewStakersMembershipValidator creates a validator for the provided
// group selection result.
func NewStakersMembershipValidator(
	stakersAddresses []relaychain.StakerAddress,
	signing chain.Signing,
) *StakersMembershipValidator {
	members := make(map[string][]int)
	for position, address := range stakersAddresses {
		addressAsString := hex.EncodeToString(address)
		positions, ok := members[addressAsString]
		if ok {
			positions = append(positions, position)
			members[addressAsString] = positions
		} else {
			members[addressAsString] = []int{position}
		}
	}

	return &StakersMembershipValidator{
		members: members,
		signing: signing,
	}
}

// IsInGroup returns true if party with the given public key has been
// selected to the group. Otherwise, function returns false.
func (smv *StakersMembershipValidator) IsInGroup(
	publicKey *operator.PublicKey,
) bool {
	address := hex.EncodeToString(
		smv.signing.PublicKeyToAddress(publicKey),
	)
	_, isInGroup := smv.members[address]
	return isInGroup
}

// IsValidMembership returns true if party with the given public key has
// been selected to the group at the given position. If the position does
// not match function returns false. The same happens when the party was
// not selected to the group.
func (smv *StakersMembershipValidator) IsValidMembership(
	memberID MemberIndex,
	publicKey []byte,
) bool {
	address := hex.EncodeToString(
		smv.signing.PublicKeyBytesToAddress(publicKey),
	)
	positions, isInGroup := smv.members[address]

	if !isInGroup {
		return false
	}

	index := int(memberID - 1)

	for _, position := range positions {
		if index == position {
			return true
		}
	}

	return false
}
