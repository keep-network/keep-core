package tbtc

import (
	"github.com/keep-network/keep-common/pkg/persistence"
	"github.com/keep-network/keep-core/pkg/ecdsa"
	"github.com/keep-network/keep-core/pkg/net"
)

// Node represents the current state of a TBTC node.
type node struct {
	chain       Chain
	netProvider net.Provider

	ecdsaNode *ecdsa.Node
}

func newNode(
	chain Chain,
	netProvider net.Provider,
	persistence persistence.Handle,
	ecdsaNode *ecdsa.Node,
) *node {
	return &node{
		chain:       chain,
		netProvider: netProvider,
		ecdsaNode:   ecdsaNode,
	}
}

func (n *node) registerNewWallet(
	walletID ecdsa.WalletID,
	walletPubKeyHash [20]byte,
) {
	// TODO: Implementation.
}
