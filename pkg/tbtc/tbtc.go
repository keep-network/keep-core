package tbtc

import (
	"context"
	"fmt"
	"github.com/ipfs/go-log"
	"github.com/keep-network/keep-common/pkg/persistence"
	"github.com/keep-network/keep-core/pkg/ecdsa"
	"github.com/keep-network/keep-core/pkg/net"
)

var logger = log.Logger("keep-tbtc")

// Initialize kicks off the TBTC by initializing internal state, ensuring
// preconditions like staking are met, and then kicking off the internal TBTC
// implementation. Returns an error if this failed.
func Initialize(
	ctx context.Context,
	chain Chain,
	netProvider net.Provider,
	persistence persistence.Handle,
) error {
	ecdsaNode, err := ecdsa.Initialize(ctx, chain, netProvider, persistence)
	if err != nil {
		return fmt.Errorf("cannot initialize ECDSA node: [%v]", err)
	}

	tbtcNode := newNode(chain, netProvider, persistence, ecdsaNode)

	_ = chain.OnNewWalletRegistered(func(event *NewWalletRegisteredEvent) {
		go func() {
			// TODO: Deduplication.

			logger.Infof(
				"New wallet with ID [0x%x] registered",
				event.WalletID,
			)

			tbtcNode.registerNewWallet(
				event.WalletID,
				event.WalletPubKeyHash,
			)
		}()
	})

	return nil
}
