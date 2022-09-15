package group

import (
	"github.com/ipfs/go-log/v2"
	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/operator"
)

// MembershipValidator operates on a group selection result and lets to
// validate one's membership based on the provided public key.
//
// Validator is used to filter out messages from parties not selected to
// the group. It is also used to confirm the position in the group of
// a party that was selected. This is used to validate messages sent by that
// party to all other group members.
type MembershipValidator struct {
	logger  log.StandardLogger
	members map[string][]int // operator address -> operator positions in group
	signing chain.Signing
}

// NewMembershipValidator creates a validator for the provided group selection
// result.
func NewMembershipValidator(
	logger log.StandardLogger,
	operatorsAddresses []chain.Address,
	signing chain.Signing,
) *MembershipValidator {
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

	return &MembershipValidator{
		logger:  logger,
		members: members,
		signing: signing,
	}
}

// IsInGroup returns true if party with the given public key has been
// selected to the group. Otherwise, function returns false.
func (mv *MembershipValidator) IsInGroup(
	publicKey *operator.PublicKey,
) bool {
	address, err := mv.signing.PublicKeyToAddress(publicKey)
	if err != nil {
		mv.logger.Errorf("cannot convert public key to chain address: [%v]", err)
		return false
	}

	_, isInGroup := mv.members[address.String()]
	return isInGroup
}

// IsValidMembership returns true if party with the given public key has
// been selected to the group at the given position. If the position does
// not match function returns false. The same happens when the party was
// not selected to the group.
func (mv *MembershipValidator) IsValidMembership(
	memberID MemberIndex,
	publicKey []byte,
) bool {
	address := mv.signing.PublicKeyBytesToAddress(publicKey).String()

	positions, isInGroup := mv.members[address]

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
