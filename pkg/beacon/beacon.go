package beacon

import (
	"context"
	"fmt"
	"math/big"
	"os"

	"github.com/keep-network/keep-core/pkg/beacon/relay"
	relaychain "github.com/keep-network/keep-core/pkg/beacon/relay/chain"
	"github.com/keep-network/keep-core/pkg/beacon/relay/event"
	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/net"
)

// Initialize kicks off the random beacon by initializing internal state,
// ensuring preconditions like staking are met, and then kicking off the
// internal random beacon implementation. Returns an error if this failed,
// otherwise enters a blocked loop.
func Initialize(
	ctx context.Context,
	stakingID string,
	relayChain relaychain.Interface,
	blockCounter chain.BlockCounter,
	stakeMonitor chain.StakeMonitor,
	netProvider net.Provider,
) error {
	chainConfig, err := relayChain.GetConfig()
	if err != nil {
		return err
	}

	staker, err := stakeMonitor.StakerFor(stakingID)
	if err != nil {
		return err
	}

	node := relay.NewNode(
		staker,
		netProvider,
		blockCounter,
		chainConfig,
	)

	// Current entry holds currently processed entry received from relay entry
	// submission. After generating a new group this entry will be used to
	//  generate a new entry. With this solution group selection and new relay
	// entry submission will be executed sequentially.
	var currentEntry *event.Entry

	relayChain.OnRelayEntryRequested(func(request *event.Request) {
		fmt.Printf("New entry requested [%+v]\n", request)

		go node.GenerateRelayEntryIfEligible(
			request.RequestID,
			request.PreviousValue,
			request.Seed,
			relayChain,
		)
	})

	relayChain.OnRelayEntryGenerated(func(entry *event.Entry) {
		fmt.Printf("Saw new relay entry [%+v]\n", entry)

		currentEntry = entry

		// new entry generated, try to join the group
		go func() {
			err := node.SubmitTicketsForGroupSelection(
				relayChain,
				blockCounter,
				entry.Value.Bytes(),
				entry.RequestID,
				entry.Seed,
			)
			if err != nil {
				fmt.Fprintf(os.Stderr, "tickets submission failed: [%v]\n", err)
			}
		}()
	})

	// TODO: This is a temporary solution until DKG Phase 14 is ready. We assume
	// that only one DKG result is published in DKG Phase 13 and submit it as a
	// final group public key.
	relayChain.OnDKGResultPublished(
		func(dkgResultPublication *event.DKGResultPublication) {
			fmt.Printf("Saw new DKG result published [%+v]\n", dkgResultPublication)

			relayChain.SubmitGroupPublicKey(
				dkgResultPublication.RequestID,
				dkgResultPublication.GroupPublicKey,
			)
		},
	)

	relayChain.OnGroupRegistered(func(registration *event.GroupRegistration) {
		fmt.Printf("Saw new group registered [%+v]\n", registration)

		node.RegisterGroup(
			registration.RequestID.String(),
			registration.GroupPublicKey,
		)

		entry := currentEntry
		nextRequestID := new(big.Int).Add(entry.RequestID, big.NewInt(1))
		go node.GenerateRelayEntryIfEligible(
			nextRequestID,
			entry.Value,
			entry.Seed,
			relayChain,
		)
	})

	return nil
}
