package ecdsa

import (
	"context"
	"fmt"
	"github.com/ipfs/go-log"
	"github.com/keep-network/keep-common/pkg/persistence"
	"github.com/keep-network/keep-core/pkg/net"
	"github.com/keep-network/keep-core/pkg/sortition"
)

var logger = log.Logger("keep-ecdsa")

// Initialize kicks off the ECDSA by initializing internal state, ensuring
// preconditions like staking are met, and then kicking off the internal ECDSA
// implementation. Returns an error if this failed, otherwise returns a handle
// to the ECDSA node.
func Initialize(
	ctx context.Context,
	chain Chain,
	netProvider net.Provider,
	persistence persistence.Handle,
) (*Node, error) {
	node := newNode(chain, netProvider, persistence)
	deduplicator := newDeduplicator()

	err := sortition.MonitorPool(ctx, chain, sortition.DefaultStatusCheckTick)
	if err != nil {
		return nil, fmt.Errorf(
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

	return node, nil
}
