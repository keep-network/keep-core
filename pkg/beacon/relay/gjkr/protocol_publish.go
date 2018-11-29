package gjkr

import (
	"math/big"

	"github.com/keep-network/keep-core/pkg/beacon/relay/event"
	"github.com/keep-network/keep-core/pkg/beacon/relay/result"
	"github.com/keep-network/keep-core/pkg/internal/sliceutils"
)

// Result returns a result of distributed key generation. It takes generated
// group public key along with disqualified and inactive members and returns
// it in a Result struct.
//
// Additional validation to check if number of disqualified and inactive members
// is greater than half of the configured dishonest threshold. If so the group
// is to weak and the result is set to a failure.
func (pm *PublishingMember) Result() *result.Result {
	group := pm.group
	disqualifiedMembers := group.DisqualifiedMembers() // DQ
	inactiveMembers := group.InactiveMembers()         // IA

	// if nPlayers(IA + DQ) > T/2:
	if len(disqualifiedMembers)+len(inactiveMembers) > (group.dishonestThreshold / 2) {
		// Result.failure(disqualified = DQ)
		return &result.Result{
			Success:      false,
			Disqualified: disqualifiedMembers,
		}
	}

	// Result.success(pubkey = Y, inactive = IA, disqualified = DQ)
	return &result.Result{
		Success:        true,
		GroupPublicKey: big.NewInt(123), // TODO: Use group public key after Phase 12 is merged
		Disqualified:   disqualifiedMembers,
		Inactive:       inactiveMembers,
	}
}

// PublishResult sends a result containing i.a. group public key to the blockchain.
// It checks if the result has already been published to the blockchain. If not
// it determines if the current member is eligable to result submission. If allowed
// it submits the results to the blockchain. The function returns result published
// to the blockchain.
//
// See Phase 13 of the protocol specification.
func (pm *PublishingMember) PublishResult(result *result.Result) (*event.PublishedResult, error) {
	chainRelay := pm.protocolConfig.ChainHandle().ThresholdRelay()

	for !chainRelay.IsResultPublished(result) { // while not resultPublished
		publishersIDs, err := pm.determinePublishersIDs()
		if err != nil {
			return nil, err
		}

		if sliceutils.Contains(publishersIDs, pm.ID) {
			errors := make(chan error)
			publishedResult := make(chan *event.PublishedResult)

			chainRelay.SubmitResult(pm.ID, result).
				OnComplete(func(pr *event.PublishedResult, err error) {
					publishedResult <- pr
					errors <- err
				})
			return <-publishedResult, <-errors
		}
	}
	return nil, nil
}
