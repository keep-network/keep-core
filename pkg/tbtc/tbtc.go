package tbtc

import (
	"context"
	"fmt"

	"github.com/ipfs/go-log"
	"github.com/keep-network/keep-common/pkg/persistence"
	"github.com/keep-network/keep-core/pkg/net"
	"github.com/keep-network/keep-core/pkg/sortition"
	"github.com/keep-network/keep-core/pkg/tecdsa"
)

// TODO: Unit tests for `tbtc.go`.

var logger = log.Logger("keep-tbtc")

// ProtocolName denotes the name of the protocol defined by this package.
const ProtocolName = "tbtc"

// Initialize kicks off the TBTC by initializing internal state, ensuring
// preconditions like staking are met, and then kicking off the internal TBTC
// implementation. Returns an error if this failed.
func Initialize(
	ctx context.Context,
	chain Chain,
	netProvider net.Provider,
	persistence persistence.Handle,
	tecdsaConfig tecdsa.Config,
) error {
	node := newNode(chain, netProvider, persistence, tecdsaConfig)
	deduplicator := newDeduplicator()

	err := sortition.MonitorPool(ctx, logger, chain, sortition.DefaultStatusCheckTick)
	if err != nil {
		return fmt.Errorf(
			"could not set up sortition pool monitoring: [%v]",
			err,
		)
	}

	_ = chain.OnDKGStarted(func(event *DKGStartedEvent) {
		go func() {
			if ok := deduplicator.notifyDKGStarted(
				event.Seed,
			); !ok {
				logger.Warningf(
					"DKG started event with seed [0x%x] and "+
						"starting block [%v] has been already processed",
					event.Seed,
					event.BlockNumber,
				)
				return
			}

			logger.Infof(
				"DKG started with seed [0x%x] at block [%v]",
				event.Seed,
				event.BlockNumber,
			)

			node.joinDKGIfEligible(
				event.Seed,
				event.BlockNumber,
			)
		}()
	})

	return nil
}
