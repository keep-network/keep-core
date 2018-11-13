package gjkr

import (
	"math/big"

	"github.com/keep-network/keep-core/pkg/beacon/relay/result"
)

// PrepareResult sets results of distributed key generation. It takes generated
// group public key along with disqualified and inactive members and stores
// in member's result field.
//
// Additional validation to check if number of disqualified and inactive members
// is greater than half of the configured dishonest threshold. If so the group
// is to weak and the result is set to a failure.
func (pm *PublishingMember) PrepareResult() {
	group := pm.group
	disqualifiedMembers := group.DisqualifiedMembers()
	inactiveMembers := group.InactiveMembers()

	// if nPlayers(IA + DQ) > T/2:
	if len(disqualifiedMembers)+len(inactiveMembers) > (group.dishonestThreshold / 2) {
		// Result.failure(disqualified = DQ)
		pm.result = &result.Result{
			Success:      false,
			Disqualified: disqualifiedMembers,
		}
	} else {
		// Result.success(pubkey = Y, inactive = IA, disqualified = DQ)
		pm.result = &result.Result{
			Success:        true,
			GroupPublicKey: big.NewInt(123), // TODO: Use group public key after Phase 12 is merged
			Disqualified:   disqualifiedMembers,
			Inactive:       inactiveMembers,
		}
	}
}
