package beacon

import (
	"context"
	"fmt"
	"os"

	"github.com/keep-network/keep-core/pkg/beacon/relay"
	relaychain "github.com/keep-network/keep-core/pkg/beacon/relay/chain"
	"github.com/keep-network/keep-core/pkg/beacon/relay/dkg"
	"github.com/keep-network/keep-core/pkg/beacon/relay/event"
	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/net"
)

type participantState int

const (
	unstaked participantState = iota
	staked
	waitingForGroup
	inIncompleteGroup
	inCompleteGroup
	inInitializingGroup
	inInitializedGroup
	inActivatingGroup
	inActiveGroup
)

// Initialize kicks off the random beacon by initializing internal state,
// ensuring preconditions like staking are met, and then kicking off the
// internal random beacon implementation. Returns an error if this failed,
// otherwise enters a blocked loop.
func Initialize(
	ctx context.Context,
	relayChain relaychain.Interface,
	blockCounter chain.BlockCounter,
	netProvider net.Provider,
) error {
	chainConfig, err := relayChain.GetConfig()
	if err != nil {
		return err
	}

	channel, err := netProvider.ChannelFor("test")
	if err != nil {
		return err
	}

	curParticipantState, err := checkParticipantState()
	if err != nil {
		panic(fmt.Sprintf("Could not resolve current relay state, aborting: [%s]", err))
	}

	switch curParticipantState {
	case unstaked:
		// check for stake command-line parameter to initialize staking?
		return fmt.Errorf("account is unstaked")
	default:
		_ = relay.NewNode(
			netProvider,
			blockCounter,
			chainConfig,
		)

		member, err := dkg.ExecuteDKG(
			blockCounter,
			channel,
			chainConfig.GroupSize,
			chainConfig.Threshold,
		)
		if err != nil {
			return err
		}

		relayChain.SubmitGroupPublicKey(
			"test",
			member.GroupPublicKeyBytes(),
		).OnSuccess(func(data *event.GroupRegistration) {
			fmt.Printf(
				"Submission of public key: [%s].\n",
				data,
			)
		}).OnFailure(func(err error) {
			fmt.Fprintf(
				os.Stderr,
				"Failed submission of public key: [%v].\n",
				err,
			)
		})

		fmt.Printf(
			"Submitting public key for member %s, group %s\n",
			member.MemberID(),
			"test",
		)
	}

	<-ctx.Done()

	return nil
}

func checkParticipantState() (participantState, error) {
	return staked, nil
}
