package maintainer

import (
	"context"
	"fmt"
	"time"

	"github.com/ipfs/go-log"
	"github.com/keep-network/keep-core/pkg/bitcoin"
	"github.com/keep-network/keep-core/pkg/chain"
)

// Back-off time which should be applied when the relay is restarted.
// It helps to avoid being flooded with error logs in case of a permanent error
// in the relay.
const restartBackoffTime = 10 * time.Second

var logger = log.Logger("keep-maintainer-relay")

// RelayChain is an interface that provides the ability to communicate with the
// relay on-chain contract.
type RelayChain interface {
	// Retarget adds a new epoch to the relay by providing a proof
	// of the difficulty before and after the retarget.
	Retarget(headers []bitcoin.BlockHeader) error

	// Ready checks whether the relay is active (i.e. genesis has been
	// performed).
	Ready() (bool, error)

	// IsAuthorizationRequired checks whether the relay requires the address
	// submitting a retarget to be authorised in advance by governance.
	IsAuthorizationRequired() (bool, error)

	// IsAuthorized checks whether the given address has been authorised to
	// submit a retarget by governance.
	IsAuthorized(address chain.Address) (bool, error)

	// Signing returns the signing associated with the chain.
	Signing() chain.Signing
}

func newRelay(
	ctx context.Context,
	btcChain bitcoin.Chain,
	chain RelayChain,
) *Relay {
	relay := &Relay{
		btcChain: btcChain,
		chain:    chain,
	}

	go relay.startControlLoop(ctx)

	return relay
}

// Relay is the part of maintainer responsible for maintaining the state of
// the relay on-chain contract.
type Relay struct {
	btcChain bitcoin.Chain
	chain    RelayChain
}

// startControlLoop launches the loop responsible for controlling the relay.
func (r *Relay) startControlLoop(ctx context.Context) {
	logger.Info("starting headers relay")

	defer func() {
		logger.Info("stopping headers relay")
	}()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			err := r.submitHeaders(ctx)
			if err != nil {
				logger.Errorf("error while submitting headers [%v]", err)
			}
		}

		time.Sleep(restartBackoffTime)
	}
}

//lint:ignore U1000 will be used soon
func (r *Relay) submitHeaders(ctx context.Context) error {
	if err := r.verifySubmissionEligibility(); err != nil {
		return fmt.Errorf(
			"relay maintainer not eligible to submit retargets [%v]",
			err,
		)
	}

	// TODO: Implement a loop that keeps submitting retargets.

	return nil
}

func (r *Relay) verifySubmissionEligibility() error {
	isReady, err := r.chain.Ready()
	if err != nil {
		return fmt.Errorf(
			"cannot check whether relay genesis has been performed [%v]",
			err,
		)
	}

	if !isReady {
		return fmt.Errorf("relay genesis has not been performed")
	}

	isAuthorizationRequired, err := r.chain.IsAuthorizationRequired()
	if err != nil {
		return fmt.Errorf(
			"cannot check whether authorization is required to submit "+
				"retargets [%v]",
			err,
		)
	}

	if !isAuthorizationRequired {
		return nil
	}

	maintainerAddress := r.chain.Signing().Address()

	isAuthorized, err := r.chain.IsAuthorized(maintainerAddress)
	if err != nil {
		return fmt.Errorf(
			"cannot check whether relay maintainer is authorized to "+
				"submit retargets [%v]",
			err,
		)
	}

	if !isAuthorized {
		return fmt.Errorf(
			"relay maintainer is not authorized to submit retargets",
		)
	}

	return nil
}
