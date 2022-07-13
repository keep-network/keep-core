package group

import (
	"github.com/keep-network/keep-core/pkg/operator"

	"github.com/keep-network/keep-core/pkg/chain"
)

// MembershipValidator lets to validate one's membership based on the
// provided public key.
type MembershipValidator interface {
	IsInGroup(publicKey *operator.PublicKey) bool
	IsValidMembership(memberID MemberIndex, publicKey []byte) bool
}

// OperatorsMembershipValidator operates on a group selection result and lets to
// validate one's membership based on the provided public key.
//
// Validator is used to filter out messages from parties not selected to
// the group. It is also used to confirm the position in the group of
// a party that was selected. This is used to validate messages sent by that
// party to all other group members.
type OperatorsMembershipValidator struct {
	members map[string][]int // operator address -> operator positions in group
	signing chain.Signing
}

// NewOperatorsMembershipValidator creates a validator for the provided
// group selection result.
func NewOperatorsMembershipValidator(
	operatorsAddresses []chain.Address,
	signing chain.Signing,
) *OperatorsMembershipValidator {
	members := make(map[string][]int)
	for position, address := range operatorsAddresses {
		addressAsString := address.String()
		positions, ok := members[addressAsString]
		if ok {
			positions = append(positions, position)
			members[addressAsString] = positions
		} else {
			members[addressAsString] = []int{position}
		}
	}

	return &OperatorsMembershipValidator{
		members: members,
		signing: signing,
	}
}

// IsInGroup returns true if party with the given public key has been
// selected to the group. Otherwise, function returns false.
func (smv *OperatorsMembershipValidator) IsInGroup(
	publicKey *operator.PublicKey,
) bool {
	address, err := smv.signing.PublicKeyToAddress(publicKey)
	if err != nil {
		logger.Errorf("cannot convert public key to chain address: [%v]", err)
		return false
	}

	_, isInGroup := smv.members[address.String()]
	return isInGroup
}

// IsValidMembership returns true if party with the given public key has
// been selected to the group at the given position. If the position does
// not match function returns false. The same happens when the party was
// not selected to the group.
func (smv *OperatorsMembershipValidator) IsValidMembership(
	memberID MemberIndex,
	publicKey []byte,
) bool {
	address := smv.signing.PublicKeyBytesToAddress(publicKey).String()

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
