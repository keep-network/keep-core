package groupselection

import (
	relaychain "github.com/keep-network/keep-core/pkg/beacon/relay/chain"
	"github.com/keep-network/keep-core/pkg/chain"
)

// Result represents the result of group selection protocol. It contains the
// list of all stakers selected to the candidate group as well as the number of
// block at which the group selection protocol completed.
type Result struct {
	SelectedStakers        []relaychain.StakerAddress
	GroupSelectionEndBlock uint64
}

// Validator is a convenience function returning MembershipValidator for the
// group selection result.
func (r *Result) Validator(signing chain.Signing) *MembershipValidator {
	return NewMembershipValidator(r.SelectedStakers, signing)
}
